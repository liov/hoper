package apidoc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/liov/hoper/go/v2/utils/log"
)

var Doc *spec.Swagger

//参数为路径和格式
func GetDoc(args ...string) *spec.Swagger {
	if Doc != nil {
		return Doc
	}
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	} else {
		return generate()
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

	if _, err := os.Stat(realPath); os.IsNotExist(err) {
		return generate()
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
			err = json.Unmarshal(data, &Doc)
			if err != nil {
				log.Error(err)
			}
		} else {
			/*var v map[string]interface{}//子类型 json: unsupported type: map[interface{}]interface{}
			//var v interface{} //json: unsupported type: map[interface{}]interface{}
			err = yaml.Unmarshal(data, &v)
			b, err := json.Marshal(&v)
			if err != nil {
				ulog.Error(err)
			}
			json.Unmarshal(b, &Doc)*/
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
					if err = json.Unmarshal(d, &Doc); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
	return Doc
}

func generate() *spec.Swagger {
	Doc = new(spec.Swagger)
	info := new(spec.Info)
	Doc.Info = info

	Doc.Swagger = "2.0"
	Doc.Paths = new(spec.Paths)
	Doc.Definitions = make(spec.Definitions)

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

	Doc.Host = "localhost:80"
	Doc.BasePath = "/"
	Doc.Schemes = []string{"http", "https"}
	Doc.Consumes = []string{"application/json"}
	Doc.Produces = []string{"application/json"}
	return Doc
}

func WriteToFile(realPath, modName string) {
	if Doc == nil {
		generate()
	}
	if realPath == "" {
		realPath = "."
	}

	realPath = realPath + modName
	err := os.MkdirAll(realPath, 0666)
	if err != nil {
		log.Error(err)
	}

	realPath = filepath.Join(realPath, modName+".service.swagger.json")

	if _, err := os.Stat(realPath); err == nil {
		os.Remove(realPath)
	}
	var file *os.File
	file, err = os.Create(realPath)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(Doc)
	if err != nil {
		log.Error(err)
	}

	/*b, err := yaml.Marshal(swag.ToDynamicJSON(Doc))
	  if err != nil {
	  	log.Error(err)
	  }
	  if _, err := file.Write(b); err != nil {
	  	log.Error(err)
	  }*/

	Doc = nil
}

func NilDoc() {
	Doc = nil
}
