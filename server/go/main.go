package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/initialize"
	"github.com/hopeio/cherry/server"
	pickgin "github.com/hopeio/pick/gin"
	chatapi "github.com/liov/hoper/server/go/chat/api"
	contentapi "github.com/liov/hoper/server/go/content/api"
	cconf "github.com/liov/hoper/server/go/content/confdao"
	uploadapi "github.com/liov/hoper/server/go/upload/api"
	upconf "github.com/liov/hoper/server/go/upload/confdao"
	userapi "github.com/liov/hoper/server/go/user/api"
	uconf "github.com/liov/hoper/server/go/user/confdao"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	// initialize.Start是提供给单服务的，这样写有问题，其他模块的配置不会更新
	defer initialize.Start(uconf.Conf, uconf.Dao)()
	defer initialize.Start(cconf.Conf, cconf.Dao)()
	defer initialize.Start(upconf.Conf, upconf.Dao)()

	config := uconf.Conf.Server.Origin()
	config.GrpcOptions = []grpc.ServerOption{grpc.StatsHandler(otelgrpc.NewServerHandler())}
	server.Start(&server.Server{
		Config: config,
		//为了可以自定义中间件
		GrpcHandler: func(gs *grpc.Server) {
			userapi.GrpcRegister(gs)
			contentapi.GrpcRegister(gs)
		},
		GinHandler: func(app *gin.Engine) {
			userapi.GinRegister(app)
			uploadapi.GinRegister(app)
			chatapi.GinRegister(app)
			contentapi.GinRegister(app)
			pickgin.Start(app, uconf.Conf.Server.GenDoc, initialize.GlobalConfig().Module, uconf.Conf.Server.Tracing)
		},
		GraphqlHandler: contentapi.NewExecutableSchema(),
		OnBeforeStart: func(ctx context.Context) {
		},
		OnAfterStart: func(ctx context.Context) {

		},
	})
}
