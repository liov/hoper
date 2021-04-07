package initialize

import (
	gini "github.com/liov/hoper/go/v2/utils/net/http/gin"
	"reflect"
	"time"
)

type ServerConfig struct {
	Protocol                               string
	Domain                                 string
	Port                                   string
	ReadTimeout                            time.Duration
	WriteTimeout                           time.Duration
	OpenTracing, SystemTracing, Prometheus bool
	Gin                                    *gini.Config
}

func (c *ServerConfig) Custom() {
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
