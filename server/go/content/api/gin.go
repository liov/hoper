package route

import (
	contentService "github.com/actliboy/hoper/server/go/content/service"
	"github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/lemon/pick"
)

func GinRegister(app *gin.Engine) {
	_ = content.RegisterMomentServiceHandlerServer(app, contentService.GetMomentService())
	_ = content.RegisterContentServiceHandlerServer(app, contentService.GetContentService())
	_ = content.RegisterActionServiceHandlerServer(app, contentService.GetActionService())
	pick.RegisterService(contentService.GetMomentService())
}
