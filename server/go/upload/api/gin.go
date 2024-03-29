package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/tiga/utils/net/http/gin/handler"
	"github.com/liov/hoper/server/go/upload/confdao"
	"github.com/liov/hoper/server/go/upload/service"
	"net/http"
)

func GinRegister(app *gin.Engine) {
	app.StaticFS("/static", http.Dir(confdao.Conf.Customize.UploadDir))
	app.GET("/api/v1/exists", handler.Convert(service.Exists))
	app.GET("/api/v1/exists/:md5/:size", service.ExistsGin)
	app.POST("/api/v1/upload/:md5", handler.Convert(service.Upload))
	app.POST("/api/v1/multiUpload", handler.Convert(service.MultiUpload))
}
