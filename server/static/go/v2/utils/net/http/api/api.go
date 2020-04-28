package api

import (
	"mime"
	"net/http"

	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

func OpenApi(mux *pick.EasyRouter, filePath string) {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	apidoc.FilePath = filePath
	mux.Handle(http.MethodGet, apidoc.PrefixUri[:len(apidoc.PrefixUri)-1], http.HandlerFunc(apidoc.ApiMod))
	mux.Handle(http.MethodGet, apidoc.PrefixUri, http.HandlerFunc(apidoc.HttpHandle))
}
