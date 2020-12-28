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
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handlerconv"
)

// 虽然我写的路由比httprouter更强大(没有map,lru cache)，但是还是选择用gin,理由是gin也用同样的方式改造了路由

func Gin(engine *gin.Engine, genApi bool,modName string) {
	methods := make(map[string]struct{})
	for _, v := range svcs {
		_, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		engine.Group(preUrl,handlerconv.Convert(middleware)...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := getMethodInfo(&method, preUrl)
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
				sess := in1.Interface().(Claims)
				sess.ParseToken(ctx.Request)
				in2 := reflect.New(in2Type.Elem())
				ctx.Bind(in2.Interface())
				result := methodValue.Call([]reflect.Value{value, in1, in2})
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
			})
			methods[methodInfo.method] = struct{}{}
		}

	}
	if genApi {
		filePath:=apidoc.FilePath
		md(filePath, modName)
		swagger(filePath, modName)
		gin_build.OpenApi(engine, filePath)
	}
	allowed := make([]string, 0, 9)
	for k := range methods {
		allowed = append(allowed, k)
	}

	registered()
}