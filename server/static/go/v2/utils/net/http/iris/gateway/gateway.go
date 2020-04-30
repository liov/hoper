package iris_gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/debug"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"
)

func Http(irisHandle func(*iris.Application), gatewayHandle gateway.GatewayHandle) http.Handler {
	gwmux := gateway.Gateway(gatewayHandle)
	//openapi
	mux := iris.New()
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		mux.Shutdown(ctx)
	})
	mux.Any("/{grpc:path}", handlerconv.FromStd(gwmux))
	mux.Any("/debug/{path:path}", handlerconv.FromStd(debug.Debug()))
	if irisHandle != nil {
		irisHandle(mux)
	}
	if err := mux.Build(); err != nil {
		log.Fatal(err)
	}
	return mux
}
