package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	"github.com/liov/hoper/go/v2/utils/net/http/api"
	"github.com/liov/hoper/go/v2/utils/net/http/debug"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

func (s *Server) Http() http.Handler {
	router := pick.NewEasyRouter(false, v2.BasicConfig.Module)
	http.DefaultServeMux.Handle("/", router)
	api.OpenApi(router, "../protobuf/api/")

	if s.GraphqlResolve != nil {
		http.DefaultServeMux.Handle("/api/graphql", handler.NewDefaultServer(s.GraphqlResolve))
	}
	if s.PickHandle != nil {
		s.PickHandle(router)
	}
	gwmux := gateway.Gateway(s.GatewayRegistr)
	//openapi
	router.Handle(pick.MethodAny, "/", gwmux)
	router.Handle(pick.MethodAny, "/debug/", debug.Debug())
	return router
}
