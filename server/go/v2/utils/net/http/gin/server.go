package gin_build

import (
	"net/http"

	"github.com/gin-gonic/gin"
	http2 "github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handlerconv"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"
)

func Http(ginHandle func(engine *gin.Engine), gatewayHandle gateway.GatewayHandle) http.Handler {

	//openapi
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	if gatewayHandle != nil {
		gwmux := gateway.Gateway(gatewayHandle)
		r.Any("/:grpc", handlerconv.FromStd(gwmux))
	}

	r.Any("/debug/:path", handlerconv.FromStd(http2.Debug()))
	if ginHandle != nil {
		ginHandle(r)
	}
	return r
}
