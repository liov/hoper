package main

import (
	"context"

	"github.com/gin-gonic/gin"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/hopeio/cherry"
	"github.com/hopeio/gox/log"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
	timex "github.com/hopeio/gox/time"
	"github.com/hopeio/pick"
	pickgin "github.com/hopeio/pick/gin"
	"github.com/hopeio/scaffold/grpc/gateway"
	"github.com/hopeio/scaffold/otel"
	commonapi "github.com/liov/hoper/server/go/common/api"
	contentapi "github.com/liov/hoper/server/go/content/api"
	uploadapi "github.com/liov/hoper/server/go/file/api"
	"github.com/liov/hoper/server/go/global"
	chatapi "github.com/liov/hoper/server/go/message/api"
	userapi "github.com/liov/hoper/server/go/user/api"
	"github.com/liov/hoper/server/go/user/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

//go:generate protogen.exe go -d -e -w -v -i ../../proto
func main() {
	defer global.Global.Cleanup()
	gatewayx.DefaultMarshal = gateway.JsonMarshal
	timex.DefaultEncoding = timex.EncodingUnixMilliseconds
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(global.Global.RootConfig.Name),
		),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}
	shutdown, err := otel.SetupOTelSDK(ctx,res)
	if err != nil {
		log.Fatalf("Failed to setup OpenTelemetry: %v", err)
	}
	global.Global.Defer(func() {
		shutdown(ctx)
	})
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
