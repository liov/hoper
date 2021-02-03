package pick

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/log"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/net/http/request/binding"
	"google.golang.org/grpc"
)

type Context interface {
	context.Context
	jwt.Claims
	grpc.ServerTransportStream
}

var (
	svcs         = make([]Service, 0)
	isRegistered = false
	claimsType   = reflect.TypeOf((*Context)(nil)).Elem()
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
		groupApiInfos = append(groupApiInfos, &groupApiInfo{describe, infos})
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

func commonHandler(w http.ResponseWriter, req *http.Request, handle *reflect.Value, ps *Params) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		for i := 0; i < handleNumIn; i++ {
			if handleTyp.In(i).Implements(claimsType) {
				params[i] = reflect.ValueOf(req.Context())
			} else {
				params[i] = reflect.New(handleTyp.In(i).Elem())
				if ps != nil || req.URL.RawQuery != "" {
					src := req.URL.Query()
					if ps != nil {
						pathParam := *ps
						if len(pathParam) > 0 {
							for i := range pathParam {
								src.Set(pathParam[i].Key, pathParam[i].Value)
							}
						}
					}
					binding.PickDecode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		resHandler(w, result)
	}
}

func resHandler(w http.ResponseWriter, result []reflect.Value) {
	if !result[1].IsNil() {
		json.NewEncoder(w).Encode(errorcode.ErrHandle(result[1].Interface()))
		return
	}
	if info, ok := result[0].Interface().(*httpi.File); ok {
		header := w.Header()
		header.Set(httpi.HeaderContentType, httpi.ContentBinaryHeaderValue)
		header.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(w, info.File)
		if flusher, canFlush := w.(http.Flusher); canFlush {
			flusher.Flush()
		}
		info.File.Close()
		return
	}
	json.NewEncoder(w).Encode(httpi.ResData{
		Code:    0,
		Message: "OK",
		Details: result[0].Interface(),
	})
}
