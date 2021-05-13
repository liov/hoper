package apidoc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/liov/hoper/v2/utils/log"
)

var PrefixUri = "/api-doc/"
var FilePath = "../protobuf/api/"

const swagger = "swagger"
const EXT = ".service.swagger.json"

func HttpHandle(w http.ResponseWriter, r *http.Request) {
	prefixUri := PrefixUri + swagger + "/"
	if r.RequestURI[len(r.RequestURI)-5:] == ".json" {
		specDoc, err := loads.Spec(FilePath + r.RequestURI[len(prefixUri):])
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		specDoc, err = specDoc.Expanded(&spec.ExpandOptions{
			SkipSchemas:         false,
			ContinueOnError:     true,
			AbsoluteCircularRef: true,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//#nosec
		_, _ = w.Write(b)
		return
	}
	mod := r.RequestURI[len(prefixUri):]
	middleware.Redoc(middleware.RedocOpts{
		BasePath: prefixUri,
		SpecURL:  path.Join(prefixUri+mod, mod+EXT),
		Path:     mod,
	}, http.NotFoundHandler()).ServeHTTP(w, r)
}

func ApiMod(w http.ResponseWriter, r *http.Request) {
	fileInfos, err := ioutil.ReadDir(FilePath)
	if err != nil {
		log.Error(err)
	}
	var ret []string
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			ret = append(ret, `<a href="`+r.RequestURI+"/"+fileInfos[i].Name()+`">`+fileInfos[i].Name()+`</a>`)
		}
	}
	w.Write([]byte(strings.Join(ret, "<br>")))
}
