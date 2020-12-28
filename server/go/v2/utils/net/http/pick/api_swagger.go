package pick

import (
	"log"
	"path/filepath"
	"reflect"

	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
)

func swagger(filePath, modName string) {
	doc := apidoc.GetDoc(filepath.Join(filePath+modName, modName+apidoc.GatewayEXT))
	for _, v := range svcs {
		describe, preUrl, _ := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := getMethodInfo(&method, preUrl)
			if methodInfo == nil {
				continue
			}
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}
			methodInfo.Swagger(doc, value.Method(j).Type(), describe, value.Type().Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}
