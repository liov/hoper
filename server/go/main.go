package main

import (
	"github.com/hopeio/pandora/pick"
	"github.com/liov/hoper/server/go/mod/chat"
	"github.com/liov/hoper/server/go/mod/upload"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/pandora/initialize"
	"github.com/hopeio/pandora/tiga"
	"github.com/hopeio/pandora/utils/log"
	cconf "github.com/liov/hoper/server/go/mod/content/confdao"
	contentService "github.com/liov/hoper/server/go/mod/content/service"
	"github.com/liov/hoper/server/go/mod/protobuf/content"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
	upconf "github.com/liov/hoper/server/go/mod/upload/confdao"
	uconf "github.com/liov/hoper/server/go/mod/user/confdao"
	userService "github.com/liov/hoper/server/go/mod/user/service"

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

	tiga.Start(&tiga.Server{
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
