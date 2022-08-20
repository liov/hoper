package server

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	gini "github.com/actliboy/hoper/server/go/lib/utils/net/http/gin"
	"reflect"
	"time"
)

type ServerConfig struct {
	Protocol                        string
	Domain                          string
	Port                            string
	ReadTimeout                     time.Duration `expr:"$+5"`
	WriteTimeout                    time.Duration `expr:"$+5"`
	OpenTracing, Prometheus, GenDoc bool
	Gin                             *gini.Config
}

func (c *ServerConfig) Init() {
	if c.Port == "" {
		c.Port = ":8080"
	}
	c.ReadTimeout = c.ReadTimeout * time.Second
	c.WriteTimeout = c.WriteTimeout * time.Second
}

func GetServiceConfig() *ServerConfig {
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
