package gin_build

import (
	"mime"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handlerconv"

	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
)

func OpenApi(mux *gin.Engine, filePath string) {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	apidoc.FilePath = filePath
	mux.GET(apidoc.PrefixUri, handlerconv.FromStd(apidoc.ApiMod))
	mux.GET(apidoc.PrefixUri+"{mod:path}", handlerconv.FromStd(apidoc.HttpHandle))
}
