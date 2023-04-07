package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/pandora/utils/net/http/gin/handler"
	"github.com/liov/hoper/server/go/mod/upload/conf"
	"github.com/liov/hoper/server/go/mod/upload/service"
	"net/http"
)

func Register(app *gin.Engine) {
	app.StaticFS("/static", http.Dir(conf.Conf.Customize.UploadDir))
	app.GET("/api/v1/exists", handler.Convert(service.Exists))
	app.GET("/api/v1/exists/:md5/:size", service.ExistsGin)
	app.POST("/api/v1/upload/:md5", handler.Convert(service.Upload))
	app.POST("/api/v1/multiUpload", handler.Convert(service.MultiUpload))
}
