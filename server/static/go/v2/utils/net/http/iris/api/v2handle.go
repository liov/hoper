package api

import (
	"io"
	"net/http"
	"reflect"

	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/response"
)

var handleMap = make(map[string]reflect.Value)

func commonHandler(ctx iris.Context) {
	handle, ok := handleMap[ctx.Path()]
	if !ok {
		ctx.NotFound()
		return
	}
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		for i := 0; i < handleNumIn; i++ {
			params[i] = reflect.New(handleTyp.In(i).Elem())
			if handleTyp.In(i).Implements(contextType) {
				sess := params[i].Interface().(Claims)
				tokenString := ctx.GetCookie("token")
				if len(tokenString) == 0 {
					tokenString = ctx.GetHeader("Authorization")
				}
				if len(tokenString) != 0 {
					sess.ParseToken(tokenString)
				}
			} else {
				//这块借鉴iris schema有待优化
				if ctx.Method() == http.MethodGet {
					queryDecoder.Decode(params[i].Interface(), ctx.Request().URL.Query())
				} else {
					ctx.ReadJSON(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		if !result[1].IsNil() {
			ctx.JSON(result[1].Interface())
			return
		}
		if info, ok := result[0].Interface().(*response.File); ok {
			ctx.Header("Content-Type", "application/octet-stream")
			ctx.Header("Content-Disposition", "attachment;filename="+info.Name)
			ctx.StreamWriter(func(w io.Writer) bool {
				io.Copy(w, info.File)
				return false
			})
			info.File.Close()
			return
		}
		ctx.JSON(response.ResData{
			Code:    0,
			Message: "success",
			Details: result[0].Interface(),
		})
	}
}

var queryDecoder = schema.NewDecoder()

func init() {
	queryDecoder.SetAliasTag("json")
}
