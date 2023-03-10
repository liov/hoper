package tiga

import (
	gini "github.com/liov/hoper/server/go/lib/utils/net/http/gin"
	"time"
)

type ServerConfig struct {
	Protocol                        string
	Domain                          string
	Port                            string
	StaticFs                        []*StaticFsConfig
	ReadTimeout                     time.Duration `expr:"$+5"`
	WriteTimeout                    time.Duration `expr:"$+5"`
	OpenTracing, Prometheus, GenDoc bool
	Gin                             *gini.Config
	GrpcWeb                         bool
	Http3                           *Http3Config
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

type Http3Config struct {
	Address  string
	CertFile string
	KeyFile  string
}

type StaticFsConfig struct {
	Prefix string
	Root   string
}
