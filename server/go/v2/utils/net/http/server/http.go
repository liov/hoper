package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/log"
	iris_build "github.com/liov/hoper/go/v2/utils/net/http/iris"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/gateway"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/middleware"
)

func (s *Server) Http() http.Handler {
	irisHandle := func(mux *iris.Application) {
		iris_build.WithConfiguration(mux, initialize.ConfUrl)
		logger := (&log.Config{Development: initialize.InitConfig.Env == initialize.PRODUCT}).NewLogger()
		middleware.SetLog(mux, logger, false)
		api.OpenApi(mux, "../protobuf/api/")

		if s.GraphqlResolve != nil {
			mux.Post("/api/graphql", handlerconv.FromStd(handler.NewDefaultServer(s.GraphqlResolve)))
		}
		if s.IrisHandle != nil {
			s.IrisHandle(mux)
		}
	}
	mux := iris_gateway.Http(irisHandle, s.GatewayRegistr)
	return mux
}
