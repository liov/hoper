package tailmon

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handlerconv"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/log"
	gin_build "github.com/liov/hoper/go/v2/utils/net/http/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/middleware"
)

func (s *Server) Http() http.Handler {
	ginHandle := func(mux *gin.Engine) {
		gin_build.WithConfiguration(mux, initialize.ConfUrl)
		logger := (&log.Config{Development: initialize.InitConfig.Env == initialize.PRODUCT}).NewLogger()
		middleware.SetLog(mux, logger, false)
		gin_build.OpenApi(mux, "../protobuf/api/")

		if s.GraphqlResolve != nil {
			mux.POST("/api/graphql", handlerconv.FromStd(handler.NewDefaultServer(s.GraphqlResolve)))
		}
		if s.GinHandle != nil {
			s.GinHandle(mux)
		}
	}
	mux := gin_build.Http(ginHandle, s.GatewayRegistr)
	return mux
}
