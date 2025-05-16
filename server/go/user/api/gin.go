package api

import (
	"github.com/gin-gonic/gin"
	model "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/service"
	"net/http"
)

func GinRegister(app *gin.Engine) {
	model.RegisterUserServiceHandlerServer(app, service.GetUserService())
	model.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
	//oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
	app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
}
