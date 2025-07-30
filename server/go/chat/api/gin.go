package api

import (
	"github.com/gin-gonic/gin"
	gini "github.com/hopeio/gox/net/http/gin"
	"github.com/liov/hoper/server/go/chat/service"
)

func GinRegister(app *gin.Engine) {
	app.GET("/api/ws/chat", gini.Convert(service.Chat))
}
