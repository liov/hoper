package gini

import (
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handler"
)

func OpenApi(mux *gin.Engine, filePath string) {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	apidoc.FilePath = filePath
	mux.GET(apidoc.PrefixUri+"md/*file", func(ctx *gin.Context) {
		mod:=ctx.Params.ByName("file")
		http.ServeFile(ctx.Writer, ctx.Request, filePath+"/"+mod+"/"+mod+"apidoc.md")
	})
	mux.GET(apidoc.PrefixUri+"swagger", handler.FromStd(apidoc.ApiMod))
	mux.GET(apidoc.PrefixUri+"swagger/*file", handler.FromStd(apidoc.HttpHandle))
}
