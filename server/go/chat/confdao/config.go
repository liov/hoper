package confdao

import (
	"github.com/hopeio/lemon/initialize/gormdb"
	"github.com/hopeio/lemon/initialize/log"
	"github.com/hopeio/lemon/initialize/redis"
	"github.com/hopeio/lemon/initialize/ristretto"
	"github.com/hopeio/lemon/initialize/server"
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
