package main

import (
	"context"

	"github.com/gin-gonic/gin"
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
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc"
	_ "github.com/hopeio/scaffold/prometheus"
)

//go:generate protogen.exe go -d -e -w -v -i ../../proto
func main() {
	defer global.Global.Cleanup()
	gatewayx.DefaultMarshal = gateway.ProtobufMarshal
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
	global.Conf.Server.WithOptions(
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
		})).Run()
}
