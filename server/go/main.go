package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry"
	"github.com/hopeio/pick"
	pickgin "github.com/hopeio/pick/gin"
	"github.com/hopeio/utils/log"
	chatapi "github.com/liov/hoper/server/go/chat/api"
	coconf "github.com/liov/hoper/server/go/common/global"
	contentapi "github.com/liov/hoper/server/go/content/api"
	cconf "github.com/liov/hoper/server/go/content/global"
	uploadapi "github.com/liov/hoper/server/go/upload/api"
	upconf "github.com/liov/hoper/server/go/upload/global"
	userapi "github.com/liov/hoper/server/go/user/api"
	uconf "github.com/liov/hoper/server/go/user/global"
	"github.com/liov/hoper/server/go/user/service"
	"google.golang.org/grpc"
)

//go:generate protogen.exe go -e -w -v -p ../../proto
func main() {
	//配置初始化应该在第一位
	// initialize.Start是提供给单服务的，这样写有问题，其他模块的配置不会更新
	defer uconf.Global.Cleanup()
	uconf.Global.Inject(cconf.Conf, cconf.Dao)
	uconf.Global.Inject(upconf.Conf, upconf.Dao)
	uconf.Global.Inject(coconf.Conf, coconf.Dao)
	log.Info("proxy:", uconf.Global.Get("proxy"))
	log.Info("proxy:", uconf.Global.RootConfig.Proxy)
	log.Info("proxy:", uconf.Global.Get("http_proxy"))
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
