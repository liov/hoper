package iris_build

import (
	"mime"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
)

func OpenApi(mux *iris.Application, filePath string) {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	apidoc.FilePath = filePath
	mux.Get(apidoc.PrefixUri, handlerconv.FromStd(apidoc.ApiMod))
	mux.Get(apidoc.PrefixUri+"{mod:path}", handlerconv.FromStd(apidoc.HttpHandle))
}
