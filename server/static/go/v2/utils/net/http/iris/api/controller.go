package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/macro"
	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

type Claims interface {
	ParseToken(string) error
}

type App iris.Application

func (a *App) Original() *iris.Application {
	return (*iris.Application)(a)
}

func ApiDoc(ctrl []Controller) {
	handler := &Handler{apiInfo: &apiInfo{}}
	for i := range ctrl {
		value := reflect.ValueOf(ctrl[i])
		if value.Kind() != reflect.Ptr {
			panic("必须传入指针")
		}

		value1 := value.Elem()
		if value1.NumField() == 0 {
			panic("controller必须有一个类型为Handler的field")
		}
		if value1.Field(0).Type() != reflect.TypeOf(handler) {
			panic("Handler field必须在第一个")
		}
		value1.Field(0).Set(reflect.ValueOf(handler))

		for j := 0; j < value.NumMethod(); j++ {
			value.Method(j).Call(nil)
		}

		path := "/api/v" + strconv.Itoa(handler.version) + handler.path
		fmt.Println(path)
	}
}

func Register(app *iris.Application, genApi bool, modName string) {
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
			if method.Type.NumOut() > 0 {
				continue
			}
			value.Method(j).Call(nil)
			if handler.path == "" || handler.method == "" || handler.title == "" ||
				handler.version == 0 || handler.createlog.version == "" {
				panic("接口路径,方法,描述,版本,创建日志均为必填")
			}

			path := "/api/v" + strconv.Itoa(handler.version) + handler.path
			var handle iris.Handler
			if handler.handle != nil {
				handle = handler.handle
			} else {
				handle = reflect.MakeFunc(reflect.TypeOf(handle),
					func(args []reflect.Value) (results []reflect.Value) {
						ctx := args[0].Interface().(iris.Context)
						ctx.WriteString("测试")
						return nil
					}).Interface().(iris.Handler)
			}
			handles := append(svcs[i].Middle(), handle)
			app.Handle(handler.method, path, handles...)

			fmt.Printf(" %s\t %s %s\t %s\n",
				pio.Green("API:"),
				pio.Yellow(strings2.FormatLen(handler.method, 6)),
				pio.Blue(strings2.FormatLen(path, 50)), pio.Purple(handler.title))
			if genApi {
				handler.Api(app)
			}
			handler.Zero()
		}
	}
	if genApi {
		apidoc.WriteToFile(apidoc.FilePath, modName)
	}
	handler = nil
}

type Controller interface {
	Middle() []iris.Handler
	Describe() string
}

type Handler struct {
	*apiInfo
	service interface{}
	handle  iris.Handler
}
type apiInfo struct {
	path, method, title string
	version             int
	changelog           []changelog
	createlog           changelog
	deprecated          *changelog
}

type changelog struct {
	version, auth, date, log string
}

var handler *Handler

func Path(p string) *Handler {
	if handler == nil {
		handler = &Handler{apiInfo: &apiInfo{}}
	}
	handler.path = p
	return handler
}

func (h *Handler) ChangeLog(v, auth, date, log string) *Handler {
	h.changelog = append(h.changelog, changelog{v, auth, date, log})
	return h
}

func (h *Handler) CreateLog(v, auth, date, log string) *Handler {
	if h.createlog.version != "" {
		panic("创建记录只允许一条")
	}
	h.createlog = changelog{v, auth, date, log}
	return h
}

func (h *Handler) Version(v int) *Handler {
	h.version = v
	return h
}

func (h *Handler) Method(m string) *Handler {
	h.method = m
	return h
}

func (h *Handler) Title(d string) *Handler {
	h.title = d
	return h
}

func (h *Handler) Service(s interface{}) *Handler {
	h.service = s
	return h
}

func (h *Handler) Handle(handle iris.Handler) *Handler {
	h.handle = handle
	return h
}

var contextType = reflect.TypeOf((*Claims)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func (h *Handler) Api(app *iris.Application) {
	doc := apidoc.GetDoc()
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
	var request, response interface{}
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
					request = reflect.New(serviceType.In(i)).Elem().Interface()
				}
			}
		}
		if numOut > 0 {
			if numOut > 2 {
				panic("service返回最多为两个")
			}
			for i := 0; i < numOut; i++ {
				if !serviceType.Out(i).Implements(errorType) {
					response = reflect.New(serviceType.Out(i)).Elem().Interface()
				}
			}
		}
	}

	if request != nil {
		idx := len(tmpl.Params)
		parameters[idx] = spec.Parameter{
			ParamProps: spec.ParamProps{
				Name: "body",
				In:   "body",
			},
		}

		parameters[idx].Schema = new(spec.Schema)
		parameters[idx].Schema.Ref = spec.MustCreateRef("#/definitions/" + reflect.TypeOf(request).Elem().Name())
		if doc.Definitions == nil {
			doc.Definitions = make(map[string]spec.Schema)
		}
		DefinitionsApi(doc.Definitions, request, exclude)
	}
	if response != nil {
		var responses spec.Responses
		responses.StatusCodeResponses = make(map[int]spec.Response)
		specResponse := spec.Response{ResponseProps: spec.ResponseProps{Schema: new(spec.Schema)}}
		specResponse.Schema.Ref = spec.MustCreateRef("#/definitions/" + reflect.TypeOf(response).Elem().Name())
		specResponse.Description = "一个成功的返回"
		DefinitionsApi(doc.Definitions, response, nil)
		responses.StatusCodeResponses[200] = specResponse
		op := spec.Operation{
			OperationProps: spec.OperationProps{
				Summary:    h.title,
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

func (h *Handler) Zero() {
	h.handle = nil
	h.deprecated = nil
	h.version = 0
	h.createlog.version = ""
	h.changelog = nil
	h.service = nil
}

type HandlerFunc func(*Handler)

func (handler HandlerFunc) Handle(hs ...iris.Handler) HandlerFunc {
	return func(h *Handler) {
	}
}
