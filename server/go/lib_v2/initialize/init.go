package initialize

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"reflect"
)

func GetConfig[T any]() *T {
	iconf := initialize.InitConfig.Config()
	value := reflect.ValueOf(iconf).Elem()
	for i := 0; i < value.NumField(); i++ {
		if conf, ok := value.Field(i).Interface().(T); ok {
			return &conf
		}
	}
	return new(T)
}
