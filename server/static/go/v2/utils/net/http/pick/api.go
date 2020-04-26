package pick

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

var isRegistered = false

type apiInfo struct {
	path, method, title string
	version             int
	changelog           []changelog
	createlog           changelog
	deprecated          *changelog
	middleware          http.HandlerFunc
}

type changelog struct {
	version, auth, date, log string
}

func Api(f func() interface{}) {
	if !isRegistered {
		panic(f())
	}
}

func Path(p string) *apiInfo {
	return &apiInfo{path: p}
}

func (api *apiInfo) Method(m string) *apiInfo {
	api.method = m
	return api
}

func (api *apiInfo) ChangeLog(v, auth, date, log string) *apiInfo {
	api.changelog = append(api.changelog, changelog{v, auth, date, log})
	return api
}

func (api *apiInfo) CreateLog(v, auth, date, log string) *apiInfo {
	if api.createlog.version != "" {
		panic("创建记录只允许一条")
	}
	api.createlog = changelog{v, auth, date, log}
	return api
}

func (api *apiInfo) Title(d string) *apiInfo {
	api.title = d
	return api
}

func (api *apiInfo) Version(v int) *apiInfo {
	api.version = v
	return api
}

func (api *apiInfo) Deprecated(v, auth, date, log string) *apiInfo {
	api.deprecated = &changelog{v, auth, date, log}
	return api
}

func (api *apiInfo) Middleware(m http.HandlerFunc) *apiInfo {
	api.middleware = m
	return api
}

//获取负责人
func (api *apiInfo) getPrincipal() string {
	if len(api.changelog) == 0 {
		return api.createlog.auth
	}
	if api.deprecated != nil {
		return api.deprecated.auth
	}
	return api.changelog[len(api.changelog)-1].auth
}

func (api *apiInfo) Api(methodType reflect.Type, tag, dec string) {
	doc := apidoc.GetDoc()
	var pathItem *spec.PathItem
	if doc.Paths != nil && doc.Paths.Paths != nil {
		if path, ok := doc.Paths.Paths[api.path]; ok {
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
	numIn := methodType.NumIn()

	if numIn == 2 {
		if !methodType.In(1).Implements(contextType) {
			if api.method == http.MethodGet {
				InType := methodType.In(1).Elem()
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
				reqName := methodType.In(1).Elem().Name()
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
				DefinitionsApi(doc.Definitions, reflect.New(methodType.In(1)).Elem().Interface(), nil)
			}
		}
	}

	if !methodType.Out(0).Implements(errorType) {
		var responses spec.Responses
		responses.StatusCodeResponses = make(map[int]spec.Response)
		response := spec.Response{ResponseProps: spec.ResponseProps{Schema: new(spec.Schema)}}
		response.Schema.Ref = spec.MustCreateRef("#/definitions/" + methodType.Out(0).Elem().Name())
		response.Description = "一个成功的返回"
		DefinitionsApi(doc.Definitions, reflect.New(methodType.Out(0)).Elem().Interface(), nil)
		responses.StatusCodeResponses[200] = response
		op := spec.Operation{
			OperationProps: spec.OperationProps{
				Summary:    api.title,
				ID:         api.path + api.method,
				Parameters: parameters,
				Responses:  &responses,
			},
		}

		var tags, desc []string
		tags = append(tags, tag, api.createlog.version)
		desc = append(desc, dec, api.createlog.log)
		for i := range api.changelog {
			tags = append(tags, api.changelog[i].version)
			desc = append(desc, api.changelog[i].log)
		}
		op.Tags = tags
		op.Description = strings.Join(desc, "\n")

		switch api.method {
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

	doc.Paths.Paths[api.path] = *pathItem
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
		json := strings.Split(body.Field(i).Tag.Get("json"), ",")[0]
		if json == "" || json == "-" {
			continue
		}
		fieldType := body.Field(i).Type
		switch fieldType.Kind() {
		case reflect.Struct:
			typ = "object"
			v = reflect.ValueOf(v).Elem().Field(i).Addr().Interface()
			subFieldName = fieldType.Name()
		case reflect.Ptr:
			typ = "object"
			v = reflect.New(fieldType.Elem()).Interface()
			subFieldName = fieldType.Elem().Name()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			typ = "integer"
		case reflect.Array, reflect.Slice:
			typ = "array"
			v = reflect.New(reflect3.GetDereferenceType(fieldType)).Interface()
			subFieldName = reflect3.GetDereferenceType(fieldType).Name()
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
		if typ == "array" {
			subSchema.Items = new(spec.SchemaOrArray)
			subSchema.Items.Schema = &spec.Schema{}
			subSchema.Items.Schema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
			DefinitionsApi(definitions, v, nil)
		}
		schema.Properties[json] = subSchema
	}
	definitions[body.Name()] = schema
}
