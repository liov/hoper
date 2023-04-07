package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/pandora/utils/net/http/gin/handler"
	"github.com/liov/hoper/server/go/mod/chat/service"
)

func Register(app *gin.Engine) {
	app.GET("/api/ws/chat", handler.Convert(service.Chat))
}
