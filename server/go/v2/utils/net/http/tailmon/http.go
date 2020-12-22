package tailmon

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handlerconv"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"

	"github.com/liov/hoper/go/v2/initialize"
	gin_build "github.com/liov/hoper/go/v2/utils/net/http/gin"
)

func (s *Server) Http() http.Handler {
	ginHandle := func(mux *gin.Engine) {
		gin_build.WithConfiguration(mux, initialize.InitConfig.ConfUrl)

		gin_build.OpenApi(mux, "../protobuf/api/")
		if s.GraphqlResolve != nil {
			mux.POST("/api/graphql", handlerconv.FromStd(handler.NewDefaultServer(s.GraphqlResolve)))
		}
		if s.GinHandle != nil {
			s.GinHandle(mux)
		}
		if s.GatewayRegistr != nil {
			gwmux := gateway.Gateway(s.GatewayRegistr)
			mux.NoRoute(handlerconv.FromStd(gwmux))// 因为httprouter的灵活性太差,一种围魏救赵的方法
		}

	}
	mux := gin_build.Http(ginHandle)
	return mux
}
