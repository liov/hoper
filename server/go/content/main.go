package main

import (
	"github.com/gin-gonic/gin"

	"github.com/hopeio/cherry/initialize"
	"github.com/hopeio/cherry/server"
	"github.com/liov/hoper/server/go/content/confdao"
	"github.com/liov/hoper/server/go/content/service"
	model "github.com/liov/hoper/server/go/protobuf/content"
	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(confdao.Conf, confdao.Dao)()

	s := server.Server{
		Config: confdao.Conf.Server.Origin(),
		GrpcHandler: func(gs *grpc.Server) {
			model.RegisterMomentServiceServer(gs, service.GetMomentService())
		},
		GinHandler: func(engine *gin.Engine) {
			_ = model.RegisterMomentServiceHandlerServer(engine, service.GetMomentService())
		},
	}
	s.Start()
}
