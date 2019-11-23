package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/kataras/iris/v12/context"
	"github.com/liov/hoper/go/v2/utils/log"
	"gopkg.in/yaml.v3"
)

func ApiMiddle(ctx context.Context) {
	currentRouteName := ctx.GetCurrentRoute().Name()[len(ctx.Method()):]

	pathItem := new(spec.PathItem)

	GetDoc("../")
	if doc.Paths !=nil && doc.Paths.Paths!=nil{
		if path, ok := doc.Paths.Paths[currentRouteName]; ok {
			pathItem = &path
		}else {
			pathItem = NewPathItem()
		}
	}else {
		doc.Paths=&spec.Paths{Paths: map[string]spec.PathItem{}}
		pathItem = NewPathItem()
	}


	parameters := make([]spec.Parameter,ctx.Params().Store.Len(),ctx.Params().Store.Len())

	params := ctx.Params().Store

	for i := range params {
		key := params[i].Key

		//val := params[i].ValueRaw
		parameters[i] = spec.Parameter{
			ParamProps:spec.ParamProps{
				Name:key,
				In:"path",
				Description:"Description",
			},
		}
	}

	if ctx.URLParam("apidoc") == "stop" {
		defer WriteToFile("../")
	}

	var res spec.Responses
	op:=spec.Operation{
		OperationProps:spec.OperationProps{
			Description:"Description",
			Consumes:[]string{"application/x-www-form-urlencoded"},
			Tags:[]string{"Tags"},
			Summary:"Summary",
			ID :"currentRouteName" + ctx.Method(),
			Parameters:parameters,
			Responses:&res,
		},
	}

	switch ctx.Method() {
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
	doc.Paths.Paths[currentRouteName] =*pathItem
	ctx.Next()
}

var doc spec.Swagger

func NewPathItem() *spec.PathItem {
	return &spec.PathItem{
		VendorExtensible: spec.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-framework": "go-swagger",
			},
		},
		PathItemProps: spec.PathItemProps{
			Parameters: [] spec.Parameter{
				{
					ParamProps:  spec.ParamProps{In: "path"},
				},
			},
		},
	}
}

func GetDoc(args ...string) {
	if doc.Swagger != "" {
		return
	}
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}
		realPath, err := filepath.Abs(targetPath)
		if err != nil {
			log.Error(err)
		}

		apiType := "json"
		if len(args) > 1 {
			apiType = args[1]
		}

		realPath = filepath.Join(realPath, "swagger."+apiType)

		if _,err:=os.Stat(realPath);os.IsNotExist(err) {
			generate()
		} else {
			file, err := os.Open(realPath)
			if err != nil {
				log.Error(err)
			}
			defer file.Close()
			data, err := ioutil.ReadAll(file)
			if err != nil {
				log.Error(err)
			}
			/*var buf bytes.Buffer
			err = json.Compact(&buf, data)
			if err != nil {
				ulog.Error(err)
			}*/
			if apiType == "json" {
				err = json.Unmarshal(data, &doc)
				if err != nil {
					log.Error(err)
				}
			} else {
				/*				var v map[string]interface{}//子类型 json: unsupported type: map[interface{}]interface{}
								//var v interface{} //json: unsupported type: map[interface{}]interface{}
								err = yaml.Unmarshal(data, &v)
								b, err := json.Marshal(&v)
								if err != nil {
									ulog.Error(err)
								}
								json.Unmarshal(b, &doc)*/
				trimmed := bytes.TrimSpace(data)
				if len(trimmed) > 0 {
					if trimmed[0] != '{' && trimmed[0] != '[' {
						yml, err := swag.BytesToYAMLDoc(trimmed)
						if err != nil {
							log.Error(err)
						}
						d, err := swag.YAMLToJSON(yml)
						if err != nil {
							log.Error(err)
						}
						if err = json.Unmarshal(d, &doc); err != nil {
							log.Error(err)
						}
					}
				}
			}
		}

}

func generate() {

	info := new(spec.Info)
	doc.Info = info

	doc.Swagger = "2.0"
	doc.Paths = new(spec.Paths)
	doc.Definitions = make(spec.Definitions)

	info.Title = "Title"
	info.Description = "Description"
	info.Version = "0.01"
	info.TermsOfService = "TermsOfService"

	var contact spec.ContactInfo
	contact.Name = "Contact Name"
	contact.Email = "Contact Mail"
	contact.URL = "Contact URL"
	info.Contact = &contact

	var license spec.License
	license.Name = "License Name"
	license.URL = "License URL"
	info.License = &license

	doc.Host = "localhost:80"
	doc.BasePath = "/"
	doc.Schemes = []string{"http", "https"}
	doc.Consumes = []string{"application/json"}
	doc.Produces = []string{"application/json"}
}

func WriteToFile(args ...string) {
	if doc.Swagger == "" {
		generate()
	}
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}
	realPath, err := filepath.Abs(targetPath)
	if err != nil {
		log.Error(err)
	}

	apiType := "json"
	if len(args) > 1 {
		apiType = args[1]
	}

	realPath = filepath.Join(realPath, "swagger."+apiType)

	if _,err:=os.Stat(realPath);err==nil {
		os.Remove(realPath)
	}
	var file *os.File
	file, err = os.Create(realPath)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()

	if apiType == "json" {
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(&doc)
		if err != nil {
			log.Error(err)
		}
	} else {
		b, err := yaml.Marshal(swag.ToDynamicJSON(&doc))
		if err != nil {
			log.Error(err)
		}
		if _, err := file.Write(b); err != nil {
			log.Error(err)
		}
	}
}
