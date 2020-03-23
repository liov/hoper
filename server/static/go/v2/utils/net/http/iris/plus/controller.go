package iris_plus

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/kataras/iris/v12/macro"
	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/net/http/api"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

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

func Register(app *iris.Application, ctrl []Controller, genApi bool) {
	handler := &Handler{apiInfo: &apiInfo{}, app: app, genApi: genApi}
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
		handler.middle = ctrl[i].Middle()
		value1.Field(0).Set(reflect.ValueOf(handler))

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Name == "Middle" {
				continue
			}
			if method.Type.NumOut() == 1 && method.Type.Out(0) == reflect.TypeOf(handler) {
				continue
			}
			value.Method(j).Call(nil)
		}
		if genApi {
			app.Get(api.PrefixUri+"{mod:path}", handlerconv.FromStd(api.HttpHandle))
			api.FilePath = "./api/"
			api.WriteToFile(api.FilePath, ctrl[i].Name())
			api.GetDoc().Paths = nil
		}
	}
	if genApi {
		api.FilePath = "./api"
		var mod []string
		for i := range ctrl {
			mod = append(mod, `<a href =">`+
				api.PrefixUri+reflect.TypeOf(ctrl[i]).Elem().Name()+`">`+
				api.PrefixUri+reflect.TypeOf(ctrl[i]).Elem().Name()+`相关接口</a>"`)
		}
		openApi := strings.Join(mod, "<br>")
		app.Get(api.PrefixUri, func(context context.Context) {
			context.WriteString(openApi)
		})
	}
	handler = nil
}

type Controller interface {
	Middle() []iris.Handler
	Name() string
}

type Handler struct {
	*apiInfo
	middle []iris.Handler
	app    *iris.Application
	genApi bool
}

type apiInfo struct {
	path, method, describe, date string
	//生成api用的参数，不优雅的实现
	request, response, service interface{}
	version                    int
	changelog                  []changelog
	createlog                  changelog
}

type changelog struct {
	version, auth, date, log string
}

func (h *Handler) Path(p string) *Handler {
	h.path = p
	return h
}

func (h *Handler) Date(d string) *Handler {
	h.date = d
	return h
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

func (h *Handler) Describe(d string) *Handler {
	h.describe = d
	return h
}

func (h *Handler) Request(r interface{}) *Handler {
	h.request = r
	return h
}

func (h *Handler) Response(r interface{}) *Handler {
	h.response = r
	return h
}

func (h *Handler) Service(s interface{}) *Handler {
	h.service = s
	return h
}

func (h *Handler) Handle(hs ...iris.Handler) *Handler {
	if h.app == nil {
		return h
	}
	if h.path == "" || h.method == "" || h.describe == "" ||
		h.version == 0 || h.request == nil || h.response == nil ||
		h.createlog.version == "" {
		panic("接口路径,方法,描述,版本,创建日志均为必填")
	}

	path := "/api/v" + strconv.Itoa(h.version) + h.path
	handles := append(h.middle, hs...)
	h.app.Handle(h.method, path, handles...)
	fmt.Printf(" %s\t %s %s\t %s\n",
		pio.Green("API:"),
		pio.Yellow(strings2.FormatLen(h.method, 6)),
		pio.Blue(strings2.FormatLen(path, 50)), pio.Purple(h.describe))
	if h.genApi {
		h.Api()
	}
	//配合CreateLog中的判断
	h.createlog.version = ""
	return h
}

func (h *Handler) Api() {
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

	tmpl, err := macro.Parse(h.path, *h.app.Macros())
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

func DefinitionsApi(definitions map[string]spec.Schema, v interface{}, exclude []string) {
	schema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:       []string{"object"},
			Properties: make(map[string]spec.Schema),
		},
	}

	body := reflect.TypeOf(v).Elem()
	var typ, subFieldName string
	for i := 0; i < body.NumField(); i++ {
		switch body.Field(i).Type.Kind() {
		case reflect.Struct:
			typ = "object"
			v = reflect.ValueOf(v).Elem().Field(i).Addr().Interface()
			subFieldName = body.Field(i).Type.Name()
		case reflect.Ptr:
			typ = "object"
			v = reflect.New(reflect.TypeOf(v).Elem().Field(i).Type.Elem()).Interface()
			subFieldName = body.Field(i).Type.Elem().Name()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			typ = "integer"
		case reflect.Array, reflect.Slice:
			typ = "array"
		case reflect.Float32, reflect.Float64:
			typ = "number"
		case reflect.String:
			typ = "string"
		case reflect.Bool:
			typ = "boolean"

		}
		subSchema := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{typ},
			},
		}
		if typ == "object" {
			subSchema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
			DefinitionsApi(definitions, v, nil)
		}
		schema.Properties[body.Field(i).Name] = subSchema
	}
	definitions[body.Name()] = schema
}

type HandlerFunc func(*Handler)

func (handler HandlerFunc) Handle(hs ...iris.Handler) HandlerFunc {
	return func(h *Handler) {
		handler(h)
		if h.app == nil {
			return
		}
		path := "/api/v" + strconv.Itoa(h.version) + h.path
		handles := append(h.middle, hs...)
		h.app.Handle(h.method, path, handles...)
		fmt.Printf(" %s\t %s %s\t %s\n",
			pio.Purple("API:"),
			pio.Yellow(strings2.FormatLen(h.method, 6)),
			pio.Blue(strings2.FormatLen(path, 50)), pio.Green(h.describe))
	}
}
