package pick

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)
func swagger(filePath, modName string) {
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
			methodInfo := getMethodInfo(value.Method(j),preUrl)
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}
			Swagger(methodInfo, value.Method(j).Type(), describe, value.Type().Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}

func Swagger(api *apiInfo,methodType reflect.Type, tag, dec string) {
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
			subType:=reflect3.GetDereferenceType(fieldType)
			subFieldName = subType.Name()
			typ = "array " + subFieldName
			switch subType.Kind() {
			case reflect.Struct,reflect.Ptr,reflect.Array, reflect.Slice:
				v = reflect.New(subType).Interface()
				typ = "array"
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
		if typ == "array"{
			subSchema.Items = new(spec.SchemaOrArray)
			subSchema.Items.Schema = &spec.Schema{}
			subSchema.Items.Schema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
			DefinitionsApi(definitions, v, nil)
		}
		schema.Properties[json] = subSchema
	}
	definitions[body.Name()] = schema
}
