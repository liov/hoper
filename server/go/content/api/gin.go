package route

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/lemon/pick"
	contentService "github.com/liov/hoper/server/go/content/service"
	"github.com/liov/hoper/server/go/protobuf/content"
)

func GinRegister(app *gin.Engine) {
	_ = content.RegisterMomentServiceHandlerServer(app, contentService.GetMomentService())
	_ = content.RegisterContentServiceHandlerServer(app, contentService.GetContentService())
	_ = content.RegisterActionServiceHandlerServer(app, contentService.GetActionService())
	pick.RegisterService(contentService.GetMomentService())
}
