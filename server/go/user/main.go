package main

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/initialize"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/user/api"
)

func main() {
	defer initialize.Start(global.Conf, global.Dao)()

	global.Conf.Server.WithOptions(cherry.WithGinHandler(api.GinRegister), cherry.WithGrpcHandler(api.GrpcRegister))
	global.Conf.Server.Run()
}
