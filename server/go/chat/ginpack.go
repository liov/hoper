package chat

import (
	"github.com/gin-gonic/gin"
)

func ChatGin(ctx *gin.Context) {
	Chat(ctx.Writer, ctx.Request)
}
