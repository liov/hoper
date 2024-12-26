package main

import (
	"github.com/hopeio/cherry"
	"github.com/liov/hoper/server/go/user/api"
)

func main() {
	server := cherry.Server{
		GrpcHandler: api.GrpcRegister,

		GinHandler: api.GinRegister,
	}
	server.Run()
}
