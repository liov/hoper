package main

import (
	"os/signal"
	"syscall"

	"github.com/liov/hoper/go/v2/gateway/internal/config"
	"github.com/liov/hoper/go/v2/gateway/internal/server"
	"github.com/liov/hoper/go/v2/initialize"
)
//go:generate protoc -I../../protobuf/ -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/gogo/protobuf/protobuf  ../../protobuf/response/*.proto --gogo_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types:../protobuf
////go:generate mv ../github.com/liov/hoper/go/v2/protobuf/response/* ../protobuf/response
////go:generate rm -r ../github.com
func main() {
	defer initialize.Start(config.Conf, nil)()
Loop:
	for {
		signal.Notify(server.SignalChan(),
			// kill -SIGINT XXXX 或 Ctrl+c
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-server.SignalChan():
			break Loop
		default:
			server.GateWay()
		}
	}
}
