package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry"
	"github.com/hopeio/initialize"
	"github.com/hopeio/initialize/conf_center/nacos"
	pickgin "github.com/hopeio/pick/gin"
	"github.com/hopeio/utils/log"
	chatapi "github.com/liov/hoper/server/go/chat/api"
	contentapi "github.com/liov/hoper/server/go/content/api"
	cconf "github.com/liov/hoper/server/go/content/confdao"
	uploadapi "github.com/liov/hoper/server/go/upload/api"
	upconf "github.com/liov/hoper/server/go/upload/confdao"
	userapi "github.com/liov/hoper/server/go/user/api"
	uconf "github.com/liov/hoper/server/go/user/confdao"
	"github.com/liov/hoper/server/go/user/service"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	// initialize.Start是提供给单服务的，这样写有问题，其他模块的配置不会更新
	defer initialize.Start(uconf.Conf, uconf.Dao, &nacos.Nacos{})()
	initialize.GlobalConfig().Inject(cconf.Conf, cconf.Dao)
	initialize.GlobalConfig().Inject(upconf.Conf, upconf.Dao)
	log.Info("proxy:", initialize.GlobalConfig().Get("proxy"))
	log.Info("proxy:", initialize.GlobalConfig().InitConfig.Proxy)
	log.Info("proxy:", initialize.GlobalConfig().Get("http_proxy"))
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
			pickgin.Register(app, uconf.Conf.Server.EnableTelemetry, &service.UserService{})
		}
		//s.GraphqlHandler= graphql.NewExecutableSchema(),
	}).Run()

}
