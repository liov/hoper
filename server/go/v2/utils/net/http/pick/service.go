package pick

import (
	"context"
	"log"
	"net/http"
	"reflect"

	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
)

type Claims interface {
	ParseToken(*http.Request) error
}

var claimsType = reflect.TypeOf((*Claims)(nil)).Elem()
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

type Service interface {
	//返回描述，url的前缀，中间件
	Service() (describe, prefix string, middleware []http.HandlerFunc)
}

var svcs = make([]Service, 0)

func RegisterService(svc ...Service) {
	svcs = append(svcs, svc...)
}

var isRegistered = false

func registered() {
	isRegistered = true
	svcs = nil
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

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Type.NumIn() < 2 || method.Type.NumOut() != 2 {
				continue
			}
			methodInfo := getMethodInfo(value.Method(j), preUrl)
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}

			router.Handle(methodInfo.method, methodInfo.path, methodInfo.middleware, value.Method(j))
			methods[methodInfo.method] = struct{}{}
			Log(methodInfo.method, methodInfo.path, describe+":"+methodInfo.title)
		}
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
