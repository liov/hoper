package server

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/tiga"
	"reflect"
)

type ServerConfig tiga.ServerConfig

func (c *ServerConfig) Init() {
	(*tiga.ServerConfig)(c).Init()
}

func (c *ServerConfig) Origin() *tiga.ServerConfig {
	return (*tiga.ServerConfig)(c)
}

func GetServerConfig() *ServerConfig {
	iconf := initialize.InitConfig.Config()
	value := reflect.ValueOf(iconf).Elem()
	for i := 0; i < value.NumField(); i++ {
		if conf, ok := value.Field(i).Interface().(ServerConfig); ok {
			return &conf
		}
	}
	return defaultServerConfig()
}

func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: ":8080",
	}
}
