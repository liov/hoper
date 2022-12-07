package tiga

import (
	gini "github.com/liov/hoper/server/go/lib/utils/net/http/gin"
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
	GrpcWeb                         bool
}

func (c *ServerConfig) Init() {
	if c.Port == "" {
		c.Port = ":8080"
	}
	c.ReadTimeout = c.ReadTimeout * time.Second
	c.WriteTimeout = c.WriteTimeout * time.Second
}

func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: ":8080",
	}
}
