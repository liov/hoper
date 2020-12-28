package pick

import (
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	gin_build "github.com/liov/hoper/go/v2/utils/net/http/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handler"
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
		engine.Group(preUrl, handler.Converts(middleware)...)
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
				in1.Interface().(Claims).ParseToken(ctx.Request)
				in2 := reflect.New(in2Type.Elem())
				ctx.Bind(in2.Interface())
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				ginResHandler(ctx,result)
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