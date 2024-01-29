package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/lemon/utils/net/http/gin/handler"
	"github.com/liov/hoper/server/go/chat/service"
)

func GinRegister(app *gin.Engine) {
	app.GET("/api/ws/chat", handler.Convert(service.Chat))
}
