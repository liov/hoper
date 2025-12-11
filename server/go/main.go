package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry"
	"github.com/hopeio/pick"
	pickgin "github.com/hopeio/pick/gin"
	_ "github.com/hopeio/scaffold/grpc/gateway"
	chatapi "github.com/liov/hoper/server/go/chat/api"
	commonapi "github.com/liov/hoper/server/go/common/api"
	contentapi "github.com/liov/hoper/server/go/content/api"
	uploadapi "github.com/liov/hoper/server/go/file/api"
	"github.com/liov/hoper/server/go/global"
	userapi "github.com/liov/hoper/server/go/user/api"
	"github.com/liov/hoper/server/go/user/service"
	"google.golang.org/grpc"
)

//go:generate protogen.exe go -d -e -w -v -i ../../proto
func main() {
	//配置初始化应该在第一位
	defer global.Global.Cleanup()
	global.Conf.Server.WithOptions(cherry.WithGinHandler(func(app *gin.Engine) {
		commonapi.GinRegister(app)
		userapi.GinRegister(app)
		uploadapi.GinRegister(app)
		chatapi.GinRegister(app)
		contentapi.GinRegister(app)
		pick.HandlerPrefix("Pick")
		pickgin.Register(app, &service.UserService{})
	}), cherry.WithGrpcHandler(func(gs *grpc.Server) {
		userapi.GrpcRegister(gs)
		contentapi.GrpcRegister(gs)
	})).Run()
}
