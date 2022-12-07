package pick

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/reflect"
	"github.com/liov/hoper/server/go/lib/utils/strings"
)

const Template = `
func (*UserService) Add(ctx *model.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() {
		pick.Post("/add").
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").End()
	})

	return &response.TinyRep{Message: req.Name}, nil
}
`

type apiInfo struct {
	path, method, title string
	version             int
	changelog           []changelog
	createlog           changelog
	deprecated          *changelog
	middleware          []http.HandlerFunc
}

type changelog struct {
	version, auth, date, log string
}

func Get(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodGet}
}
func Post(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodPost}
}
func Put(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodPut}
}
func Delete(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodDelete}
}
func Patch(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodPatch}
}
func Trace(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodTrace}
}
func Head(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodHead}
}
func Options(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodOptions}
}
func Connect(p string) *apiInfo {
	return &apiInfo{path: p, method: http.MethodConnect}
}

func (api *apiInfo) ChangeLog(v, auth, date, log string) *apiInfo {
	v = version(v)
	api.changelog = append(api.changelog, changelog{v, auth, date, log})
	return api
}

func version(v string) string {
	if v[0] != 'v' {
		return "v" + v
	}
	return v
}

func (api *apiInfo) CreateLog(v, auth, date, log string) *apiInfo {
	if api.createlog.version != "" {
		panic("创建记录只允许一条")
	}
	v = version(v)
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
	v = version(v)
	api.deprecated = &changelog{v, auth, date, log}
	return api
}

func (api *apiInfo) Middleware(m ...http.HandlerFunc) *apiInfo {
	api.middleware = m
	return api
}

func (api *apiInfo) End() {
	panic(api)
}

// 获取负责人
func (api *apiInfo) getPrincipal() string {
	if len(api.changelog) == 0 {
		return api.createlog.auth
	}
	if api.deprecated != nil {
		return api.deprecated.auth
	}
	return api.changelog[len(api.changelog)-1].auth
}

// 简直就是精髓所在，真的是脑洞大开才能想到
func getMethodInfo(method *reflect.Method, preUrl string, claimsTyp reflect.Type) (info *apiInfo) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo); ok {
				//_,_, info.version = parseMethodName(method.Name)
				if v.version == 0 {
					v.version = 1
				}
				v.path = preUrl + "/v" + strconv.Itoa(v.version) + v.path
				info = v
			} else {
				log.Error(err)
			}
		}
	}()
	methodValue := method.Func
	methodType := methodValue.Type()
	numIn := methodType.NumIn()
	numOut := methodType.NumOut()
	var err error
	defer func() {
		if err != nil {
			log.Debugf("%s %s 未注册:%v", preUrl, method.Name, err)
		}
	}()

	if numIn != 3 {
		err = errors.New("method参数必须为两个")
		return
	}
	if numOut != 2 {
		err = errors.New("method返回值必须为两个")
		return
	}
	if !methodType.In(1).Implements(claimsTyp) {
		return
	}
	if !methodType.Out(1).Implements(errorType) {
		err = errors.New("service第二个返回值必须为error类型")
		return
	}
	params := make([]reflect.Value, numIn, numIn)
	for i := 0; i < numIn; i++ {
		params[i] = reflect.New(methodType.In(i).Elem())
	}
	methodValue.Call(params)
	return nil
}

// 从方法名称分析出接口名和版本号
func parseMethodName(originName string, methods []string) (method, name string, version int) {
	idx := strings.LastIndexByte(originName, 'V')
	version = 1
	if idx > 0 {
		if v, err := strconv.Atoi(originName[idx+1:]); err == nil {
			version = v
		} else {
			idx = len(originName)
		}
	} else {
		idx = len(originName)
	}
	name = stringsi.LowerFirst(originName[:idx])
	for _, method := range methods {
		if strings.HasPrefix(name, method) {
			return method, name[len(method):], version
		}
	}
	return http.MethodPost, name, version
}

func (api *apiInfo) Swagger(doc *spec.Swagger, methodType reflect.Type, tag, dec string) {
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
		if !methodType.In(1).Implements(claimsType) {
			if api.method == http.MethodGet {
				InType := methodType.In(1).Elem()
				for j := 0; j < InType.NumField(); j++ {
					param := spec.Parameter{
						ParamProps: spec.ParamProps{
							Name: InType.Field(j).Name,
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
	var arraySubType string
	for i := 0; i < body.NumField(); i++ {
		json := strings.Split(body.Field(i).Tag.Get("json"), ",")[0]
		if json == "" || json == "-" {
			continue
		}
		fieldType := body.Field(i).Type
		switch fieldType.Kind() {
		case reflect.Struct:
			typ = "object"
			v = reflect.New(fieldType).Interface()
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
			subType := reflecti.GetDereferenceType(fieldType)
			subFieldName = subType.Name()
			switch subType.Kind() {
			case reflect.Struct, reflect.Ptr, reflect.Array, reflect.Slice:
				v = reflect.New(subType).Interface()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				arraySubType = "integer"
			case reflect.Float32, reflect.Float64:
				arraySubType = "number"
			case reflect.String:
				arraySubType = "string"
			case reflect.Bool:
				arraySubType = "boolean"
			}
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
			if arraySubType == "" {
				subSchema.Items.Schema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
				DefinitionsApi(definitions, v, nil)
			} else {
				subSchema.Items.Schema.Type = []string{arraySubType}
			}

		}
		schema.Properties[json] = subSchema
	}
	definitions[body.Name()] = schema
}
