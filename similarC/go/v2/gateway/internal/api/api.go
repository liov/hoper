package api

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/liov/hoper/go/v2/utils/log"
)

var apiPath = "/open-api/"

func OpenApi(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")
	handler := middleware.Spec(apiPath, b, http.NotFoundHandler())
	mux.Handle(apiPath, http.StripPrefix(apiPath, http.FileServer(http.Dir("./api"))))
}

func httpHandle(w http.ResponseWriter, r *http.Request) {
	mod := r.RequestURI[len(apiPath):]
	specDoc, err := loads.Spec("../protobuf/" + mod + "/" + mod + ".service.swagger.json")
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

	if r.URL.Query().Get("f") == "redoc" {
		handler := middleware.Redoc(middleware.RedocOpts{
			BasePath: apiPath,
			SpecURL:  path.Join(apiPath, "swagger.json"),
			Path:     mod,
		}, http.NotFoundHandler())

	} else {
		u, err := url.Parse("http://petstore.swagger.io/")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		q := u.Query()
		q.Add("url", fmt.Sprintf("http://hoper.xyz%s", path.Join(apiPath, "swagger.json")))
		u.RawQuery = q.Encode()
	}
}
