package api

import (
	"github.com/actliboy/hoper/server/go/chat/service"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/tailmon/utils/net/http/gin/handler"
)

func GinRegister(app *gin.Engine) {
	app.GET("/api/ws/chat", handler.Convert(service.Chat))
}
