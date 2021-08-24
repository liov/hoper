package gini

import (
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/server/go/lib/utils/net/http/api/apidoc"
	"github.com/liov/hoper/server/go/lib/utils/net/http/gin/handler"
)

func OpenApi(mux *gin.Engine, filePath string) {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	apidoc.FilePath = filePath
	mux.GET(apidoc.PrefixUri+"md/*file", func(ctx *gin.Context) {
		mod := ctx.Params.ByName("file")
		http.ServeFile(ctx.Writer, ctx.Request, filePath+"/"+mod+"/"+mod+"apidoc.md")
	})
	mux.GET(apidoc.PrefixUri+"swagger", handler.Wrap(apidoc.ApiMod))
	mux.GET(apidoc.PrefixUri+"swagger/*file", handler.Wrap(apidoc.HttpHandle))
}
