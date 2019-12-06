package api

import (
	"encoding/json"
	"mime"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
)

var apiPath = "/open-api/"

func OpenApi(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	mux.Handle(apiPath, http.HandlerFunc(httpHandle))
}

func httpHandle(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI[len(r.RequestURI)-5:] == ".json" {
		specDoc, err := loads.Spec("../protobuf/" + r.RequestURI[len(apiPath):])
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
	mod := r.RequestURI[len(apiPath):]
	middleware.Redoc(middleware.RedocOpts{
		BasePath: apiPath,
		SpecURL:  path.Join(apiPath+mod, mod+".service.swagger.json"),
		Path:     mod,
	}, http.NotFoundHandler()).ServeHTTP(w, r)
}
