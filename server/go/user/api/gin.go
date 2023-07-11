package api

import (
	model "github.com/actliboy/hoper/server/go/protobuf/user"
	"github.com/actliboy/hoper/server/go/user/service"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/dora/pick"
	"net/http"
)

func GinRegister(app *gin.Engine) {
	_ = model.RegisterUserServiceHandlerServer(app, service.GetUserService())
	_ = model.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
	//oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
	app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
	pick.RegisterService(service.GetUserService())
}
