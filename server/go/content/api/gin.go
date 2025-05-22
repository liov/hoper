package route

import (
	"github.com/gin-gonic/gin"
	contentService "github.com/liov/hoper/server/go/content/service"
	"github.com/liov/hoper/server/go/protobuf/content"
)

func GinRegister(app *gin.Engine) {
	content.RegisterMomentServiceHandlerServer(app, contentService.GetMomentService())
	content.RegisterContentServiceHandlerServer(app, contentService.GetContentService())
	content.RegisterActionServiceHandlerServer(app, contentService.GetActionService())
}
