package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry"
	"github.com/hopeio/pick"
	pickgin "github.com/hopeio/pick/gin"
	chatapi "github.com/liov/hoper/server/go/chat/api"
	contentapi "github.com/liov/hoper/server/go/content/api"
	"github.com/liov/hoper/server/go/global"
	uploadapi "github.com/liov/hoper/server/go/upload/api"
	userapi "github.com/liov/hoper/server/go/user/api"
	uconf "github.com/liov/hoper/server/go/user/global"
	"github.com/liov/hoper/server/go/user/service"
	"google.golang.org/grpc"
)

//go:generate protogen.exe go -e -w -v -p ../../proto
func main() {
	//配置初始化应该在第一位
	defer global.Global.Cleanup()
	uconf.Conf.Server.WithOptions(func(s *cherry.Server) {
		s.GrpcHandler = func(gs *grpc.Server) {
			userapi.GrpcRegister(gs)
			contentapi.GrpcRegister(gs)
		}
		s.GinHandler = func(app *gin.Engine) {
			userapi.GinRegister(app)
			uploadapi.GinRegister(app)
			chatapi.GinRegister(app)
			contentapi.GinRegister(app)
			pick.HandlerPrefix("Pick")
			pickgin.Register(app, &service.UserService{})
		}
		//s.GraphqlHandler= graphql.NewExecutableSchema(),
	}).Run()

}
