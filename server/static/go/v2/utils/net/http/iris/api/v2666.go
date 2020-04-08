package api

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/kataras/iris/v12"
	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

var isRegistered = false

func Api(f func() interface{}) {
	if !isRegistered {
		panic(f())
	}
}

func Method(m string) *apiInfo {
	return &apiInfo{method: m}
}

func (h *apiInfo) Date(d string) *apiInfo {
	h.date = d
	return h
}

func (h *apiInfo) ChangeLog(v, auth, date, log string) *apiInfo {
	h.changelog = append(h.changelog, changelog{v, auth, date, log})
	return h
}

func (h *apiInfo) CreateLog(v, auth, date, log string) *apiInfo {
	if h.createlog.version != "" {
		panic("创建记录只允许一条")
	}
	h.createlog = changelog{v, auth, date, log}
	return h
}

func (h *apiInfo) Describe(d string) *apiInfo {
	h.describe = d
	return h
}

func (h *apiInfo) Auth(a string) *apiInfo {
	h.auth = a
	return h
}

type Service interface {
	Describe() string
	Middle() []iris.Handler
}

var svcs = make(map[string]Service)

func RegisterService(svc Service, tag string) {
	svcs[tag] = svc
}

func registered() {
	isRegistered = true
	svcs = nil
}

func RegisterAllService(app *iris.Application, genApi bool, modName string) {
	for k, v := range svcs {
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			panic("必须传入指针")
		}

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Name == "Middle" || method.Type.NumOut() == 1 {
				continue
			}

			methodInfo := getMethodInfo(value.Method(j))
			if methodInfo == nil {
				log.Fatalf("%s未注册", method.Name)
			}
			if methodInfo.method == "" || methodInfo.describe == "" || methodInfo.createlog.version == "" {
				log.Fatalf("接口路径,方法,描述,创建日志均为必填")
			}
			mName, version := parseMethodName(method.Name)
			methodInfo.path = "/api/v" + strconv.Itoa(version) + "/" + k + "/" + mName
			handles := append(v.Middle(), commonHandler)
			app.Handle(methodInfo.method, methodInfo.path, handles...)
			handleMap[methodInfo.path] = value.Method(j)
			fmt.Printf(" %s\t %s %s\t %s\n",
				pio.Green("API:"),
				pio.Yellow(strings2.FormatLen(methodInfo.method, 6)),
				pio.Blue(strings2.FormatLen(methodInfo.path, 50)), pio.Purple(methodInfo.describe))
			if genApi {
				methodInfo.Api(value.Method(j).Type())
				apidoc.WriteToFile(apidoc.FilePath, modName)
			}
		}
	}
	app.Get("/api-doc/html", HTML)
	app.Get("/api-doc/md", MD)

	registered()
}

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
	serviceType := fv.Type()
	params := make([]reflect.Value, 0, fv.Type().NumIn())
	numIn := serviceType.NumIn()
	numOut := serviceType.NumOut()
	if numIn == 0 {
		panic("service至少一个参数且参数必须实现Claims接口")
	}
	if numIn > 2 {
		panic("service参数最多为两个")
	}
	if numOut != 2 {
		panic("service返回值必须为两个")
	}
	if !serviceType.In(0).Implements(contextType) {
		panic("service第一个参数必须实现Claims接口")
	}
	if !serviceType.Out(1).Implements(errorType) {
		panic("service第二个返回值必须为error类型")
	}
	for i := 0; i < numIn; i++ {
		params = append(params, reflect.New(serviceType.In(i).Elem()))
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
	name = lowerFirst(originName[:idx])
	return
}

// 仅首位小写（更符合接口的规范）
func lowerFirst(t string) string {
	b := []byte(t)
	if 'A' <= b[0] && b[0] <= 'Z' {
		b[0] += 'a' - 'A'
	}
	return string(b)
}

func (h *apiInfo) Api(serviceType reflect.Type) {
	doc := apidoc.GetDoc()
	var pathItem *spec.PathItem
	if doc.Paths != nil && doc.Paths.Paths != nil {
		if path, ok := doc.Paths.Paths[h.path]; ok {
			pathItem = &path
		} else {
			pathItem = new(spec.PathItem)
		}
	} else {
		doc.Paths = &spec.Paths{Paths: map[string]spec.PathItem{}}
		pathItem = new(spec.PathItem)
	}

	//我觉得路径参数并没有那么值得非用不可
	parameters := make([]spec.Parameter, 0)
	numIn := serviceType.NumIn()

	if numIn == 2 {
		if !serviceType.In(1).Implements(contextType) {
			if h.method == http.MethodGet {
				InType := serviceType.In(1).Elem()
				for j := 0; j < InType.NumField(); j++ {
					param := spec.Parameter{
						ParamProps: spec.ParamProps{
							Name: InType.Field(1).Name,
							In:   "query",
						},
					}
					parameters = append(parameters, param)
				}
			} else {
				reqName := serviceType.In(1).Elem().Name()
				param := spec.Parameter{
					ParamProps: spec.ParamProps{
						Name: reqName,
						In:   "body",
					},
				}

				param.Schema = new(spec.Schema)
				param.Schema.Ref = spec.MustCreateRef("#/definitions/" + reqName)
				parameters = append(parameters, param)
				if doc.Definitions == nil {
					doc.Definitions = make(map[string]spec.Schema)
				}
				DefinitionsApi(doc.Definitions, reflect.New(serviceType.In(1)).Elem().Interface(), nil)
			}
		}
	}

	if !serviceType.Out(0).Implements(errorType) {
		var responses spec.Responses
		responses.StatusCodeResponses = make(map[int]spec.Response)
		response := spec.Response{ResponseProps: spec.ResponseProps{Schema: new(spec.Schema)}}
		response.Schema.Ref = spec.MustCreateRef("#/definitions/" + serviceType.Out(0).Elem().Name())
		response.Description = "一个成功的返回"
		DefinitionsApi(doc.Definitions, reflect.New(serviceType.Out(0)).Elem().Interface(), nil)
		responses.StatusCodeResponses[200] = response
		op := spec.Operation{
			OperationProps: spec.OperationProps{
				Summary:    h.describe,
				ID:         h.path + h.method,
				Parameters: parameters,
				Responses:  &responses,
			},
		}

		var tags, desc []string
		tags = append(tags, h.createlog.version)
		desc = append(desc, h.createlog.log)
		for i := range h.changelog {
			tags = append(tags, h.changelog[i].version)
			desc = append(desc, h.changelog[i].log)
		}
		op.Tags = tags
		op.Description = strings.Join(desc, "\n")

		switch h.method {
		case http.MethodGet:
			pathItem.Get = &op
		case http.MethodPost:
			pathItem.Post = &op
		case http.MethodPut:
			pathItem.Put = &op
		case http.MethodDelete:
			pathItem.Delete = &op
		case http.MethodOptions:
			pathItem.Options = &op
		case http.MethodPatch:
			pathItem.Patch = &op
		case http.MethodHead:
			pathItem.Head = &op
		}
	}

	doc.Paths.Paths[h.path] = *pathItem
}
