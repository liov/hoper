package pick

import (
	"path/filepath"

	"github.com/actliboy/hoper/server/go/lib/utils/net/http/api/apidoc"
)

func swagger(filePath, modName string) {
	doc := apidoc.GetDoc(filepath.Join(filePath+modName, modName+apidoc.EXT))
	for _, groupApiInfo := range groupApiInfos {
		for _, methodInfo := range groupApiInfo.infos {
			methodInfo.Swagger(doc, methodInfo.method, groupApiInfo.describe, methodInfo.method.Name())
		}
	}
	apidoc.WriteToFile(filePath, modName)
}
