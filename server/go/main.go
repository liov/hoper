package main

import (
	pickgin "github.com/hopeio/pick/gin"
	"github.com/hopeio/tiga/server"
	chatapi "github.com/liov/hoper/server/go/chat/api"
	contentapi "github.com/liov/hoper/server/go/content/api"
	uploadapi "github.com/liov/hoper/server/go/upload/api"
	userapi "github.com/liov/hoper/server/go/user/api"
	"go.opencensus.io/zpages"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/tiga/initialize"
	"github.com/hopeio/tiga/utils/log"
	cconf "github.com/liov/hoper/server/go/content/confdao"
	upconf "github.com/liov/hoper/server/go/upload/confdao"
	uconf "github.com/liov/hoper/server/go/user/confdao"

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
	zpages.Handle(nil, "/")
	view.RegisterExporter(&exporter.PrintExporter{})
	view.SetReportingPeriod(time.Second)
	// GinRegister the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	config := uconf.Conf.Server.Origin()
	config.GRPCOptions = []grpc.ServerOption{grpc.StatsHandler(&ocgrpc.ServerHandler{})}
	server.Start(&server.Server{
		Config: config,
		//为了可以自定义中间件

		GRPCHandle: func(gs *grpc.Server) {
			userapi.GrpcRegister(gs)
			contentapi.GrpcRegister(gs)
		},
		GinHandle: func(app *gin.Engine) {
			userapi.GinRegister(app)
			uploadapi.GinRegister(app)
			chatapi.GinRegister(app)
			contentapi.GinRegister(app)
			pickgin.Start(app, uconf.Conf.Server.GenDoc, initialize.GlobalConfig.Module, uconf.Conf.Server.Trace)
		},
		GraphqlHandle: contentapi.NewExecutableSchema(),
	})
}
