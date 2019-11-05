package main

////go:generate protoc -I ../../../protobuf/ ../../../protobuf/user/*.proto --go_out=plugins=grpc:../protobuf
//go:generate protoc -I../../protobuf/ -I$GOPATH/src -I$GOPATH/src/github.com/gogo/protobuf/protobuf  ../../protobuf/user/*.proto --gogo_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types:../protobuf
import (
	"os"
	"os/signal"
	"syscall"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/user/internal/server"
)

func main() {
	defer initialize.Start(config.Conf, dao.Dao)()
	ch := make(chan os.Signal, 1)
Loop:
	for {
		signal.Notify(ch,
			// kill -SIGINT XXXX 或 Ctrl+c
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			break Loop
		default:
			server.Server()
		}
	}
}