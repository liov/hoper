package pick

import (
	"context"
	"net/http"
	"reflect"

	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
)

type Claims interface {
	ParseToken(*http.Request,string) error
}

var (
	svcs         = make([]Service, 0)
	isRegistered = false
	claimsType   = reflect.TypeOf((*Claims)(nil)).Elem()
	contextType  = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType    = reflect.TypeOf((*error)(nil)).Elem()
)

type Service interface {
	//返回描述，url的前缀，中间件
	Service() (describe, prefix string, middleware []http.HandlerFunc)
}

func RegisterService(svc ...Service) {
	svcs = append(svcs, svc...)
}

func registered() {
	isRegistered = true
	svcs = nil
	groupApiInfos = nil
}

func Api(f func() interface{}) {
	if !isRegistered {
		panic(f())
	}
}

func register(router *Router, genApi bool, modName string) {
	methods := make(map[string]struct{})
	for _, v := range svcs {
		describe, preUrl, middleware := v.Service()
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

			router.Handle(methodInfo.method, methodInfo.path, methodInfo.middleware, value.Method(j))
			methods[methodInfo.method] = struct{}{}
			Log(methodInfo.method, methodInfo.path, describe+":"+methodInfo.title)
			infos = append(infos, &apiDocInfo{methodInfo, method.Type})
		}
		groupApiInfos = append(groupApiInfos, &groupApiInfo{
			describe: describe,
			infos:    infos,
		})
		router.GroupUse(preUrl, middleware...)
	}
	if genApi {
		OpenApi(router, apidoc.FilePath, modName)
	}
	allowed := make([]string, 0, 9)
	for k := range methods {
		allowed = append(allowed, k)
	}
	router.globalAllowed = allowedMethod(allowed)

	registered()
}
