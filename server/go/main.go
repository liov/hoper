package main

import (
	"context"

	"github.com/gin-gonic/gin"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/hopeio/cherry"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
	timex "github.com/hopeio/gox/time"
	"github.com/hopeio/pick"
	pickgin "github.com/hopeio/pick/gin"
	"github.com/hopeio/scaffold/grpc/gateway"
	commonapi "github.com/liov/hoper/server/go/common/api"
	contentapi "github.com/liov/hoper/server/go/content/api"
	uploadapi "github.com/liov/hoper/server/go/file/api"
	"github.com/liov/hoper/server/go/global"
	chatapi "github.com/liov/hoper/server/go/message/api"
	userapi "github.com/liov/hoper/server/go/user/api"
	"github.com/liov/hoper/server/go/user/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

//go:generate protogen.exe go -d -e -w -v -i ../../proto
func main() {
	//配置初始化应该在第一位
	defer global.Global.Cleanup()
	gatewayx.DefaultMarshal = gateway.JsonMarshal
	timex.DefaultEncoding = timex.EncodingUnixMilliseconds
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	prometheus.MustRegister(srvMetrics)
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	global.Conf.Server.WithOptions(
		cherry.WithGrpc(func(option *cherry.GrpcConfig) {
			option.UnaryServerInterceptors = []grpc.UnaryServerInterceptor{
				srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			}
			option.StreamServerInterceptors = []grpc.StreamServerInterceptor{
				srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			}
		}),
		cherry.WithTelemetry(func(telemetry *cherry.TelemetryConfig) {
			telemetry.StdoutExportOpts = []stdoutmetric.Option{}
		}),
		cherry.WithGinHandler(func(app *gin.Engine) {
			commonapi.GinRegister(app)
			userapi.GinRegister(app)
			uploadapi.GinRegister(app)
			chatapi.GinRegister(app)
			contentapi.GinRegister(app)
			pick.HandlerPrefix("Pick")
			pickgin.Register(app, &service.UserService{})
		}), cherry.WithGrpcHandler(func(gs *grpc.Server) {
			userapi.GrpcRegister(gs)
			contentapi.GrpcRegister(gs)
			srvMetrics.InitializeMetrics(gs)
		})).Run()
}
