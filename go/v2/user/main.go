package main

//go:generate protoc -I ../../../protobuf/ ../../../protobuf/user/*.proto --go_out=plugins=grpc:../protobuf
import (
	"flag"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/protobuf/user"
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
	log.Info(reflect.TypeOf(time.Now()).Size())
	log.Info(user.User{CreatedAt:&timestamp.Timestamp{Seconds:1,Nanos:1}})
	server.Server()
}
