package confdao

import (
	"github.com/hopeio/tailmon/initialize/gormdb"
	"github.com/hopeio/tailmon/initialize/log"
	"github.com/hopeio/tailmon/initialize/redis"
	"github.com/hopeio/tailmon/initialize/ristretto"
	"github.com/hopeio/tailmon/initialize/server"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.ServerConfig
	GORMDB    gormdb.DatabaseConfig
	Redis     redis.Config
	Cache     ristretto.CacheConfig
	Log       log.LogConfig
}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

}

type serverConfig struct {
	Host string
}
