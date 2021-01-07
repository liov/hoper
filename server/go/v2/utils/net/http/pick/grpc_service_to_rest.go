package pick

import (
	"context"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handler"
)

type MapRouter map[string]methodHandle

type AuthCtx func(r *http.Request) context.Context
// Deprecated:这种方法不推荐使用了，目前就两种定义api的方式，一种grpc-gateway，一种pick自定义
// 该方法适用于不使用grpc-gateway的情况，只用该方法定义api
func GrpcServiceToRestfulApi(engine *gin.Engine, authCtx AuthCtx, genApi bool, modName string) {
	httpMethods := []string{http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodConnect, http.MethodHead, http.MethodTrace}
	doc := apidoc.GetDoc(filepath.Join(apidoc.FilePath+modName,modName+apidoc.GatewayEXT))
	methods := make(map[string]struct{})
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		group := engine.Group(preUrl, handler.Converts(middleware)...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodType := method.Type
			methodValue := method.Func
			if method.Type.NumIn() < 3 || method.Type.NumOut() != 2 ||
				!methodType.In(1).Implements(contextType) ||
				!methodType.Out(1).Implements(errorType){
				continue
			}

			methodInfo:=new(apiInfo)
			methodInfo.title = describe
			methodInfo.middleware = middleware
			methodInfo.method, methodInfo.path, methodInfo.version = parseGrpcMethodName(method.Name, httpMethods)
			methodInfo.path = "v" + strconv.Itoa(methodInfo.version) + "/" + methodInfo.path

			in2Type := methodType.In(2)
			group.Handle(methodInfo.method, methodInfo.path, func(ctx *gin.Context) {
				in0 := reflect.ValueOf(authCtx(ctx.Request))
				in2 := reflect.New(in2Type.Elem())
				ctx.Bind(in2.Interface())
				result := methodValue.Call([]reflect.Value{value, in0, in2})
				ginResHandler(ctx,result)
			})
			methods[methodInfo.method] = struct{}{}
			if genApi {
				methodInfo.Swagger(doc,value.Method(j).Type(), describe, value.Type().Name())
			}
		}

	}
	if genApi {
		apidoc.WriteToFile(apidoc.FilePath, modName)
	}

}

func parseGrpcMethodName(name string, methods []string) (string, string, int) {
	name, version := parseMethodName(name)
	for _, method := range methods {
		if strings.HasPrefix(name, method) {
			return method, name[len(method):], version
		}
	}
	return http.MethodPost, name, version
}
