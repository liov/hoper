package main

//go:generate protoc -I ../../../protobuf/ ../../../protobuf/user/*.proto --go_out=plugins=grpc:../protobuf
//go:generate protoc -I../../protobuf/ -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf  ../../protobuf/user/*.proto --gogo_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types:../protobuf
import (
	"flag"
	"reflect"
	"time"

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
	log.Info(reflect.TypeOf(time.Now()).Size())
	server.Server()
}
