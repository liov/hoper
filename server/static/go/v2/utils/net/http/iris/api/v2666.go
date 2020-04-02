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
	"github.com/liov/hoper/go/v2/utils/net/http/api"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

var registered = false

func Api(f func() interface{}) {
	if !registered {
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
	Name() string
	DocTag() string
	Middle() []iris.Handler
}

func RegisterAllService(app *iris.Application, svcs []Service, genApi bool, modName string) {
	for i := range svcs {
		value := reflect.ValueOf(svcs[i])
		if value.Kind() != reflect.Ptr {
			panic("必须传入指针")
		}

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Name == "Middle" {
				continue
			}
			if method.Type.NumOut() == 1 {
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
			methodInfo.path = "/api/v" + strconv.Itoa(version) + "/" + svcs[i].Name() + "/" + mName
			handles := append(svcs[i].Middle(), commonHandler)
			app.Handle(methodInfo.method, methodInfo.path, handles...)
			fmt.Printf(" %s\t %s %s\t %s\n",
				pio.Green("API:"),
				pio.Yellow(strings2.FormatLen(methodInfo.method, 6)),
				pio.Blue(strings2.FormatLen(methodInfo.path, 50)), pio.Purple(methodInfo.describe))
			if genApi {
				methodInfo.Api(value.Method(j).Type())
				api.WriteToFile(api.FilePath, modName)
			}
		}
	}
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
	if numIn > 0 {
		if numIn > 2 {
			panic("service参数最多为两个")
		}
		for i := 0; i < numIn; i++ {
			params = append(params, reflect.New(serviceType.In(i).Elem()))
		}
		fv.Call(params)
	} else {
		fv.Call(nil)
	}

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
	doc := api.GetDoc()
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
	numOut := serviceType.NumOut()
	if numIn > 0 {
		if numIn > 2 {
			panic("service参数最多为两个")
		}
		for i := 0; i < numIn; i++ {
			if !serviceType.In(i).Implements(contextType) {
				if h.method == http.MethodGet {
					InType := serviceType.In(i).Elem()
					for j := 0; j < InType.NumField(); j++ {
						param := spec.Parameter{
							ParamProps: spec.ParamProps{
								Name: InType.Field(i).Name,
								In:   "query",
							},
						}
						parameters = append(parameters, param)
					}
				} else {
					reqName := serviceType.In(i).Elem().Name()
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
					DefinitionsApi(doc.Definitions, reflect.New(serviceType.In(i)).Elem().Interface(), nil)
				}
			}
		}
	}
	if numOut > 0 {
		if numOut > 2 {
			panic("service返回最多为两个")
		}
		for i := 0; i < numOut; i++ {
			if !serviceType.Out(i).Implements(errorType) {
				var responses spec.Responses
				responses.StatusCodeResponses = make(map[int]spec.Response)
				response := spec.Response{ResponseProps: spec.ResponseProps{Schema: new(spec.Schema)}}
				response.Schema.Ref = spec.MustCreateRef("#/definitions/" + serviceType.Out(i).Elem().Name())
				response.Description = "一个成功的返回"
				DefinitionsApi(doc.Definitions, reflect.New(serviceType.Out(i)).Elem().Interface(), nil)
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
		}
	}

	doc.Paths.Paths[h.path] = *pathItem
}
