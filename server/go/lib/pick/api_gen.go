package pick

import (
	contexti "github.com/actliboy/hoper/server/go/lib/tiga/context"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/api/apidoc"
	gin_build "github.com/actliboy/hoper/server/go/lib/utils/net/http/gin"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/gin/handler"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
)

func GenApi(middlewareHandler func(preUrl string, middleware []http.HandlerFunc), handle func(method, path string, in2Type reflect.Type, methodValue, value reflect.Value)) {
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		middlewareHandler(preUrl, middleware)
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*apiDocInfo

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := getMethodInfo(&method, preUrl, claimsType)
			if methodInfo == nil {
				continue
			}
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}
			methodType := method.Type
			methodValue := method.Func
			in2Type := methodType.In(2)
			handle(methodInfo.method, methodInfo.path, in2Type, methodValue, value)
			infos = append(infos, &apiDocInfo{methodInfo, method.Type})
		}
		groupApiInfos = append(groupApiInfos, &groupApiInfo{describe, infos})
	}

	registered()
}

func genApiDoc(modName string) {
	filePath := apidoc.FilePath
	md(filePath, modName)
	swagger(filePath, modName)
}

func GenGinAPI(genApi bool, modName string, engine *gin.Engine) {
	GenApi(func(preUrl string, middleware []http.HandlerFunc) {
		engine.Group(preUrl, handler.Converts(middleware)...)
	},
		func(method, path string, in2Type reflect.Type, methodValue, value reflect.Value) {
			engine.Handle(method, path, func(ctx *gin.Context) {
				ctxi, span := contexti.CtxFromRequest(ctx.Request, true)
				if span != nil {
					defer span.End()
				}
				in1 := reflect.ValueOf(ctxi)
				in2 := reflect.New(in2Type.Elem())
				gin_build.Bind(ctx, in2.Interface())
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				resHandler(ctxi, ctx.Writer, result)
			})
		})
	if genApi {
		genApiDoc(modName)
		gin_build.OpenApi(engine, apidoc.FilePath)
	}

}
