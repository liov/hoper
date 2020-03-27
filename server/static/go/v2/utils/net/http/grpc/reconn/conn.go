package reconn

import (
	"reflect"
	"unsafe"

	"github.com/liov/hoper/go/v2/initialize/v2"
	"google.golang.org/grpc"
)

func ReConnect(v interface{}, module string, opts []grpc.DialOption) func() error {
	value := reflect.ValueOf(v).Elem()
	ptr := value.Field(0).Pointer()
	conn := (*grpc.ClientConn)(unsafe.Pointer(ptr))
	return func() error {
		conn.Close()
		newConn, err := grpc.Dial(initialize.BasicConfig.NacosConfig.GetServiceEndPort(module), opts...)
		if err != nil {
			return err
		}
		newConnPtr := (**grpc.ClientConn)(unsafe.Pointer(value.Addr().Pointer()))
		*newConnPtr = newConn

		return nil
	}
}
