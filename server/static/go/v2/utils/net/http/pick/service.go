package pick

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

type Service interface {
	//返回描述，url的前缀，中间件
	Service() (describe, prefix string, middleware http.HandlerFunc)
}

var svcs = make([]Service, 0)

func RegisterService(svc ...Service) {
	svcs = append(svcs, svc...)
}

func registered() {
	isRegistered = true
	svcs = nil
}

func New(genApi bool, modName string) *Router {
	router := &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		middleware:             make([]http.HandlerFunc, 0),
	}
	methods := make(map[string]struct{})
	for _, v := range svcs {
		describe, preUrl, _ := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Type.NumIn() < 2 || method.Type.NumOut() != 2 {
				continue
			}
			methodInfo := getMethodInfo(value.Method(j))
			if methodInfo == nil {
				log.Fatalf("%s未注册", method.Name)
			}
			if methodInfo.version == 0 {
				methodInfo.version = 1
			}
			methodInfo.path = preUrl + methodInfo.path
			methodInfo.path = strings.Replace(methodInfo.path, "${version}", "v"+strconv.Itoa(methodInfo.version), 1)
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}

			router.Handle(methodInfo.method, methodInfo.path, value.Method(j))
			methods[methodInfo.method] = struct{}{}
			fmt.Printf(" %s\t %s %s\t %s\n",
				pio.Green("API:"),
				pio.Yellow(strings2.FormatLen(methodInfo.method, 6)),
				pio.Blue(strings2.FormatLen(methodInfo.path, 50)), pio.Purple(methodInfo.title))
			if genApi {
				methodInfo.Api(value.Method(j).Type(), describe, value.Type().Name())
				apidoc.WriteToFile(apidoc.FilePath, modName)
			}
		}
	}

	allowed := make([]string, 0, 9)
	for k := range methods {
		allowed = append(allowed, k)
	}
	router.globalAllowed = allowedMethod(allowed)

	registered()
	return router
}

type Claims interface {
	ParseToken(*http.Request) error
}

var contextType = reflect.TypeOf((*Claims)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

//简直就是精髓所在，真的是脑洞大开才能想到
func getMethodInfo(fv reflect.Value) (info *apiInfo) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo); ok {
				info = v
			} else {
				log.Panic(err)
			}
		}
	}()
	methodType := fv.Type()
	params := make([]reflect.Value, 0, fv.Type().NumIn())
	numIn := methodType.NumIn()
	numOut := methodType.NumOut()
	if numIn == 1 {
		panic("method至少一个参数且参数必须实现Claims接口")
	}
	if numIn > 2 {
		panic("method参数最多为两个")
	}
	if numOut != 2 {
		panic("method返回值必须为两个")
	}
	if !methodType.In(0).Implements(contextType) {
		panic("service第一个参数必须实现Claims接口")
	}
	if !methodType.Out(1).Implements(errorType) {
		panic("service第二个返回值必须为error类型")
	}
	for i := 0; i < numIn; i++ {
		params = append(params, reflect.New(methodType.In(i).Elem()))
	}
	fv.Call(params)
	return nil
}

// 从方法名称分析出接口名和版本号
func parseMethodName(originName string) (name string, version int) {
	idx := strings.LastIndexByte(originName, 'V')
	version = 1
	if idx > 0 {
		if v, err := strconv.Atoi(originName[idx+1:]); err == nil {
			version = v
		}
	} else {
		idx = len(originName)
	}
	name = strings2.LowerFirst(originName[:idx])
	return
}
