package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cconf "github.com/liov/hoper/go/v2/content/conf"
	cdao "github.com/liov/hoper/go/v2/content/dao"
	contentervice "github.com/liov/hoper/go/v2/content/service"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	uconf "github.com/liov/hoper/go/v2/user/conf"
	udao "github.com/liov/hoper/go/v2/user/dao"
	userservice "github.com/liov/hoper/go/v2/user/service"
	"github.com/liov/hoper/go/v2/utils/log"
	grpci "github.com/liov/hoper/go/v2/utils/net/http/grpc"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/liov/hoper/go/v2/utils/net/http/tailmon"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(uconf.Conf, udao.Dao)()
	defer initialize.Start(cconf.Conf, cdao.Dao)()
	view.RegisterExporter(&exporter.PrintExporter{})
	view.SetReportingPeriod(time.Second)
	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}
	pick.RegisterService(userservice.GetUserService(), contentervice.GetMomentService())
	(&tailmon.Server{
		GRPCOptions: []grpc.ServerOption{grpci.DefaultUnaryInterceptor()},
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
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			pick.Gin(app, user.ConvertContext, true, initialize.InitConfig.Module)
		},
		CustomContext:  user.CtxWithRequest,
		ConvertContext: user.ConvertContext,
	}).Start()
}
