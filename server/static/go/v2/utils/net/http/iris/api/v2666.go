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
	"github.com/kataras/iris/v12/macro"
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

func Path(p string) *apiInfo {
	return &apiInfo{path: p}
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

func (h *apiInfo) Version(v int) *apiInfo {
	h.version = v
	return h
}

func (h *apiInfo) Method(m string) *apiInfo {
	h.method = m
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

func RegisterAllService(app *iris.Application, svcs []Service, genApi bool) {
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
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.describe == "" ||
				methodInfo.version == 0 || methodInfo.createlog.version == "" {
				log.Fatalf("接口路径,方法,描述,版本,创建日志均为必填")
			}
			path := "/api/v" + strconv.Itoa(methodInfo.version) + "/" + svcs[i].Name() + methodInfo.path
			handles := append(svcs[i].Middle(), commonHandler)
			app.Handle(methodInfo.method, path, handles...)
			fmt.Printf(" %s\t %s %s\t %s\n",
				pio.Green("API:"),
				pio.Yellow(strings2.FormatLen(methodInfo.method, 6)),
				pio.Blue(strings2.FormatLen(path, 50)), pio.Purple(h.describe))
			if genApi {
				methodInfo.Api()
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

func (h *apiInfo) Api(app *iris.Application) {
	doc := api.GetDoc()
	var pathItem *spec.PathItem
	path := "/api/v" + strconv.Itoa(h.version) + h.path
	if doc.Paths != nil && doc.Paths.Paths != nil {
		if path, ok := doc.Paths.Paths[path]; ok {
			pathItem = &path
		} else {
			pathItem = new(spec.PathItem)
		}
	} else {
		doc.Paths = &spec.Paths{Paths: map[string]spec.PathItem{}}
		pathItem = new(spec.PathItem)
	}

	tmpl, err := macro.Parse(h.path, *app.Macros())
	if err != nil {
		panic(err)
	}
	var exclude = make([]string, len(tmpl.Params))
	parameters := make([]spec.Parameter, len(tmpl.Params)+1)
	for i := 0; i < len(tmpl.Params); i++ {
		parameters[i] = spec.Parameter{
			ParamProps: spec.ParamProps{
				Name:     tmpl.Params[i].Name,
				In:       "path",
				Required: true,
			},
		}
		parameters[i].Type = "string"
		parameters[i].Format = tmpl.Params[i].Type.Indent()
		exclude[i] = tmpl.Params[i].Name
	}

	if h.service != nil {
		serviceType := reflect.TypeOf(h.service)
		numIn := serviceType.NumIn()
		numOut := serviceType.NumOut()
		if numIn > 0 {
			if numIn > 2 {
				panic("service参数最多为两个")
			}
			for i := 0; i < numIn; i++ {
				if !serviceType.In(i).Implements(contextType) {
					h.request = reflect.New(serviceType.In(i)).Interface()
				}
			}
		}
		if numOut > 0 {
			if numOut > 2 {
				panic("service返回最多为两个")
			}
			for i := 0; i < numOut; i++ {
				if !serviceType.Out(i).Implements(errorType) {
					h.response = reflect.New(serviceType.Out(i)).Interface()
				}
			}
		}
	}

	if h.request != nil {
		idx := len(tmpl.Params)
		parameters[idx] = spec.Parameter{
			ParamProps: spec.ParamProps{
				Name: "body",
				In:   "body",
			},
		}

		parameters[idx].Schema = new(spec.Schema)
		parameters[idx].Schema.Ref = spec.MustCreateRef("#/definitions/" + reflect.TypeOf(h.request).Elem().Name())
		if doc.Definitions == nil {
			doc.Definitions = make(map[string]spec.Schema)
		}
		DefinitionsApi(doc.Definitions, h.request, exclude)
	}
	if h.response != nil {
		var responses spec.Responses
		responses.StatusCodeResponses = make(map[int]spec.Response)
		response := spec.Response{ResponseProps: spec.ResponseProps{Schema: new(spec.Schema)}}
		response.Schema.Ref = spec.MustCreateRef("#/definitions/" + reflect.TypeOf(h.response).Elem().Name())
		response.Description = "一个成功的返回"
		DefinitionsApi(doc.Definitions, h.response, nil)
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
	doc.Paths.Paths[path] = *pathItem
}
