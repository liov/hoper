package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginx "github.com/hopeio/gox/net/http/gin"
	"github.com/hopeio/scaffold/gin/wrap"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/file/service"
)

func GinRegister(app *gin.Engine) {
	app.StaticFS("/upload", http.Dir(global.Conf.Customize.UploadDir))
	app.StaticFS("/static", http.Dir("D:\\Download"))
	app.GET("/api/preUpload", wrap.HandlerWrapGRPC(service.GetFileService().PreUpload))
	app.GET("/api/preUpload/:md5/:size", wrap.HandlerWrapGRPC(service.GetFileService().PreUpload))
	app.POST("/api/upload/:md5", ginx.Convert(service.Upload))
	app.POST("/api/multiUpload", ginx.Convert(service.MultiUpload))
}
