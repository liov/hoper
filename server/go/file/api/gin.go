package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/scaffold/gin/warp"
	gini "github.com/hopeio/utils/net/http/gin"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/file/service"
	"net/http"
)

func GinRegister(app *gin.Engine) {
	app.StaticFS("/upload", http.Dir(global.Conf.Customize.UploadDir))
	app.StaticFS("/static", http.Dir("D:\\Download"))
	app.GET("/api/v1/exists", warp.HandlerWrapCompatibleGRPC(service.GetFileService().Exists))
	app.GET("/api/v1/exists/:md5/:size", warp.HandlerWrapCompatibleGRPC(service.GetFileService().Exists))
	app.POST("/api/v1/upload/:md5", gini.Convert(service.Upload))
	app.POST("/api/v1/multiUpload", gini.Convert(service.MultiUpload))
}
