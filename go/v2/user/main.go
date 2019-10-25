package main

//go:generate protoc -I ../../../protobuf/ ../../../protobuf/user/*.proto --go_out=plugins=grpc:../protobuf
import (
	"flag"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/user/internal/server"
	"github.com/liov/hoper/go/v2/utils/log"
)

func main() {
	flag.Parse()
	defer log.Sync()
	initialize.Start(config.Conf,dao.Dao)
	defer dao.Dao.Close()
	server.Server()
}
