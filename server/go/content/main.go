package main

import (
	"github.com/gin-gonic/gin"

	"github.com/hopeio/initialize"
	"github.com/liov/hoper/server/go/content/global"
	"github.com/liov/hoper/server/go/content/service"
	model "github.com/liov/hoper/server/go/protobuf/content"
	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(global.Conf, global.Dao)()

	s := server.Server{
		Config: global.Conf.Server.Origin(),
		GrpcHandler: func(gs *grpc.Server) {
			model.RegisterMomentServiceServer(gs, service.GetMomentService())
		},
		GinHandler: func(engine *gin.Engine) {
			_ = model.RegisterMomentServiceHandlerServer(engine, service.GetMomentService())
		},
	}
	s.Start()
}
