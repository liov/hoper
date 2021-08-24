package initialize

import (
	gini "github.com/liov/hoper/server/go/lib/utils/net/http/gin"
	"reflect"
	"time"
)

type ServerConfig struct {
	Protocol                               string
	Domain                                 string
	Port                                   string
	ReadTimeout                            time.Duration `expr:"$+5"`
	WriteTimeout                           time.Duration `expr:"$+5"`
	OpenTracing, SystemTracing, Prometheus bool
	Gin                                    *gini.Config
}

func (c *ServerConfig) Init() {
	if c.Port == "" {
		c.Port = ":8080"
	}
	c.ReadTimeout = c.ReadTimeout * time.Second
	c.WriteTimeout = c.WriteTimeout * time.Second
}

func (init *Init) GetServiceConfig() *ServerConfig {
	value := reflect.ValueOf(init.conf).Elem()
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
