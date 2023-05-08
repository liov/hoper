package main

import (
	"github.com/hopeio/pandora/pick"
	"github.com/hopeio/pandora/server"
	"github.com/liov/hoper/server/go/chat"
	"github.com/liov/hoper/server/go/upload"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/pandora/initialize"
	"github.com/hopeio/pandora/utils/log"
	cconf "github.com/liov/hoper/server/go/content/confdao"
	contentService "github.com/liov/hoper/server/go/content/service"
	"github.com/liov/hoper/server/go/protobuf/content"
	"github.com/liov/hoper/server/go/protobuf/user"
	upconf "github.com/liov/hoper/server/go/upload/confdao"
	uconf "github.com/liov/hoper/server/go/user/confdao"
	userService "github.com/liov/hoper/server/go/user/service"

	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	// initialize.Start是提供给单服务的，这样写有问题，其他模块的配置不会更新
	defer initialize.Start(uconf.Conf, uconf.Dao)()
	defer initialize.Start(cconf.Conf, cconf.Dao)()
	defer initialize.Start(upconf.Conf, upconf.Dao)()
	view.RegisterExporter(&exporter.PrintExporter{})
	view.SetReportingPeriod(time.Second)
	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	server.Start(&server.Server{
		Config: uconf.Conf.Server.Origin(),
		//为了可以自定义中间件
		GRPCOptions: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(),
			grpc.ChainStreamInterceptor(),
			//grpc.StatsHandler(&ocgrpc.ServerHandler{})
		},
		GRPCHandle: func(gs *grpc.Server) {
			user.RegisterUserServiceServer(gs, userService.GetUserService())
			user.RegisterOauthServiceServer(gs, userService.GetOauthService())
			content.RegisterMomentServiceServer(gs, contentService.GetMomentService())
			content.RegisterContentServiceServer(gs, contentService.GetContentService())
			content.RegisterActionServiceServer(gs, contentService.GetActionService())
		},
		GinHandle: func(app *gin.Engine) {
			_ = user.RegisterUserServiceHandlerServer(app, userService.GetUserService())
			_ = user.RegisterOauthServiceHandlerServer(app, userService.GetOauthService())
			_ = content.RegisterMomentServiceHandlerServer(app, contentService.GetMomentService())
			_ = content.RegisterContentServiceHandlerServer(app, contentService.GetContentService())
			_ = content.RegisterActionServiceHandlerServer(app, contentService.GetActionService())
			app.Static("/oauth/login", "./static/login.html")
			upload.Register(app)
			chat.Register(app)
			pick.RegisterService(userService.GetUserService(), contentService.GetMomentService())
			pick.Gin(app, uconf.Conf.Server.GenDoc, initialize.InitConfig.Module, uconf.Conf.Server.OpenTracing)
		},
	})
}
