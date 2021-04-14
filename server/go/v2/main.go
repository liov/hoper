package main

import (
	"github.com/liov/hoper/go/v2/tailmon/pick"
	"github.com/liov/hoper/go/v2/upload"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handler"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cconf "github.com/liov/hoper/go/v2/content/conf"
	cdao "github.com/liov/hoper/go/v2/content/dao"
	contentervice "github.com/liov/hoper/go/v2/content/service"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/tailmon"
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	upconf "github.com/liov/hoper/go/v2/upload/conf"
	updao "github.com/liov/hoper/go/v2/upload/dao"
	uconf "github.com/liov/hoper/go/v2/user/conf"
	udao "github.com/liov/hoper/go/v2/user/dao"
	userservice "github.com/liov/hoper/go/v2/user/service"
	"github.com/liov/hoper/go/v2/utils/log"

	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(uconf.Conf, udao.Dao)()
	defer initialize.Start(cconf.Conf, cdao.Dao)()
	defer initialize.Start(upconf.Conf, updao.Dao)()
	view.RegisterExporter(&exporter.PrintExporter{})
	view.SetReportingPeriod(time.Second)
	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}
	pick.RegisterService(userservice.GetUserService(), contentervice.GetMomentService())
	(&tailmon.Server{
		GRPCOptions: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(),
			grpc.ChainStreamInterceptor(),
			//grpc.StatsHandler(&ocgrpc.ServerHandler{})
		},
		//为了可以自定义中间件
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
			app.Static("/static", "F:/upload")
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			app.GET("/api/v1/exists", handler.Convert(upload.Exists))
			app.GET("/api/v1/exists/:md5/:size", upload.ExistsGin)
			app.POST("/api/v1/upload/:md5", handler.Convert(upload.Upload))
			app.POST("/api/v1/multiUpload", handler.Convert(upload.MultiUpload))
			pick.Gin(app, true, initialize.InitConfig.Module)
		},
	}).Start()
}
