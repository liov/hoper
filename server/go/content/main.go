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

	global.Conf.Server.GrpcHandler = func(gs *grpc.Server) {
		model.RegisterMomentServiceServer(gs, service.GetMomentService())
	}
	global.Conf.Server.GinHandler = func(engine *gin.Engine) {
		model.RegisterMomentServiceHandlerServer(engine, service.GetMomentService())
	}
	global.Conf.Server.Run()
}
