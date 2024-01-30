package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/tiga/pick"
	model "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/service"
	"net/http"
)

func GinRegister(app *gin.Engine) {
	_ = model.RegisterUserServiceHandlerServer(app, service.GetUserService())
	_ = model.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
	//oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
	app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
	pick.RegisterService(service.GetUserService())
}
