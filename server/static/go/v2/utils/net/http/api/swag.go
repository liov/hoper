package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/liov/hoper/go/v2/utils/log"
)

var PrefixUri = "/api-doc/"
var FilePath = "../protobuf/"

func HttpHandle(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI[len(r.RequestURI)-5:] == ".json" {
		specDoc, err := loads.Spec(FilePath + r.RequestURI[len(PrefixUri):])
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		specDoc, err = specDoc.Expanded(&spec.ExpandOptions{
			SkipSchemas:         false,
			ContinueOnError:     true,
			AbsoluteCircularRef: true,
		})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//#nosec
		_, _ = w.Write(b)
		return
	}
	mod := r.RequestURI[len(PrefixUri):]
	middleware.Redoc(middleware.RedocOpts{
		BasePath: PrefixUri,
		SpecURL:  path.Join(PrefixUri+mod, mod+".service.swagger.json"),
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
			ret = append(ret, `<a href="`+PrefixUri+fileInfos[i].Name()+`">`+fileInfos[i].Name()+`</a>`)
		}
	}
	w.Write([]byte(strings.Join(ret, "<br>")))
}
