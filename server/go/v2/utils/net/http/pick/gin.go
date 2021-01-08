package pick

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	gin_build "github.com/liov/hoper/go/v2/utils/net/http/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handler"
)

// 虽然我写的路由比httprouter更强大(没有map,lru cache)，但是还是选择用gin,理由是gin也用同样的方式改造了路由

func Gin(engine *gin.Engine, genApi bool,modName string) {
	for _, v := range svcs {
		_, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		engine.Group(preUrl, handler.Converts(middleware)...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := getMethodInfo(&method, preUrl,claimsType)
			if methodInfo == nil{
				continue
			}
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}
			methodType := method.Type
			methodValue := method.Func
			in2Type := methodType.In(2)
			engine.Handle(methodInfo.method, methodInfo.path, func(ctx *gin.Context) {
				in1 := reflect.New(methodType.In(1).Elem())
				in1.Interface().(Claims).ParseToken(ctx.Request)
				in2 := reflect.New(in2Type.Elem())
				gin_build.Bind(ctx,in2.Interface())
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				ginResHandler(ctx,result)
			})
		}

	}
	if genApi {
		filePath:=apidoc.FilePath
		md(filePath, modName)
		swagger(filePath, modName)
		gin_build.OpenApi(engine, filePath)
	}
	registered()
}

func ginResHandler(ctx *gin.Context,result []reflect.Value)  {
	if !result[1].IsNil() {
		json.NewEncoder(ctx.Writer).Encode(result[1].Interface())
		return
	}
	if info, ok := result[0].Interface().(*httpi.File); ok {
		header := ctx.Writer.Header()
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Disposition", "attachment;filename="+info.Name)
		io.Copy(ctx.Writer, info.File)
		if flusher, canFlush := ctx.Writer.(http.Flusher); canFlush {
			flusher.Flush()
		}
		info.File.Close()
		return
	}
	ctx.JSON(200, result[0].Interface())
}