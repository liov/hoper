package main

import (
	"github.com/actliboy/hoper/server/go/lib/pick"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/gin/handler"
	"github.com/actliboy/hoper/server/go/mod/chat"
	"github.com/actliboy/hoper/server/go/mod/upload"
	"net/http"
	"time"

	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	cconf "github.com/actliboy/hoper/server/go/mod/content/conf"
	cdao "github.com/actliboy/hoper/server/go/mod/content/dao"
	contentervice "github.com/actliboy/hoper/server/go/mod/content/service"
	"github.com/actliboy/hoper/server/go/mod/protobuf/content"
	"github.com/actliboy/hoper/server/go/mod/protobuf/user"
	upconf "github.com/actliboy/hoper/server/go/mod/upload/conf"
	updao "github.com/actliboy/hoper/server/go/mod/upload/dao"
	uconf "github.com/actliboy/hoper/server/go/mod/user/conf"
	udao "github.com/actliboy/hoper/server/go/mod/user/dao"
	userservice "github.com/actliboy/hoper/server/go/mod/user/service"
	"github.com/gin-gonic/gin"

	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	// initialize.Start是提供给单服务的，这样写有问题，其他模块的配置不会更新
	defer initialize.Start(uconf.Conf, udao.Dao)()
	defer initialize.Start(cconf.Conf, cdao.Dao)()
	defer initialize.Start(upconf.Conf, updao.Dao)()
	view.RegisterExporter(&exporter.PrintExporter{})
	view.SetReportingPeriod(time.Second)
	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	(&tiga.Server{
		//为了可以自定义中间件
		GRPCOptions: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(),
			grpc.ChainStreamInterceptor(),
			//grpc.StatsHandler(&ocgrpc.ServerHandler{})
		},
		GRPCHandle: func(gs *grpc.Server) {
			user.RegisterUserServiceServer(gs, userservice.GetUserService())
			user.RegisterOauthServiceServer(gs, userservice.GetOauthService())
			content.RegisterMomentServiceServer(gs, contentervice.GetMomentService())
			content.RegisterContentServiceServer(gs, contentervice.GetContentService())
			content.RegisterActionServiceServer(gs, contentervice.GetActionService())
		},
		GinHandle: func(app *gin.Engine) {
			_ = user.RegisterUserServiceHandlerServer(app, userservice.GetUserService())
			_ = user.RegisterOauthServiceHandlerServer(app, userservice.GetOauthService())
			_ = content.RegisterMomentServiceHandlerServer(app, contentervice.GetMomentService())
			_ = content.RegisterContentServiceHandlerServer(app, contentervice.GetContentService())
			_ = content.RegisterActionServiceHandlerServer(app, contentervice.GetActionService())
			app.Static("/static", string(upconf.Conf.Customize.UploadDir))
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			app.GET("/api/v1/exists", handler.Convert(upload.Exists))
			app.GET("/api/v1/exists/:md5/:size", upload.ExistsGin)
			app.POST("/api/v1/upload/:md5", handler.Convert(upload.Upload))
			app.POST("/api/v1/multiUpload", handler.Convert(upload.MultiUpload))
			app.GET("/api/ws/chat", handler.Convert(chat.Chat))
			pick.RegisterService(userservice.GetUserService(), contentervice.GetMomentService())
			pick.Gin(app, uconf.Conf.Server.GenDoc, initialize.InitConfig.Module, uconf.Conf.Server.OpenTracing)
		},
	}).Start()
}
