package api

import (
	"github.com/gin-gonic/gin"
	gini "github.com/hopeio/utils/net/http/gin"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/file/service"
	"net/http"
)

func GinRegister(app *gin.Engine) {
	app.StaticFS("/static", http.Dir(global.Conf.Customize.UploadDir))
	app.GET("/api/v1/exists", gini.Convert(service.Exists))
	app.GET("/api/v1/exists/:md5/:size", service.ExistsGin)
	app.POST("/api/v1/upload/:md5", gini.Convert(service.Upload))
	app.POST("/api/v1/multiUpload", gini.Convert(service.MultiUpload))
}
