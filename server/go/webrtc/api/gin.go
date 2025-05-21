package api

import (
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/server/go/webrtc/service"
)

func GinRegister(app *gin.Engine) {
	app.GET("/video/*file", service.Video)
	app.GET("/live/stream", service.Stream)
}
