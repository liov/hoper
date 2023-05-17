package chat

import (
	"github.com/actliboy/hoper/server/go/chat/service"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/pandora/utils/net/http/gin/handler"
)

func Register(app *gin.Engine) {
	app.GET("/api/ws/chat", handler.Convert(service.Chat))
}
