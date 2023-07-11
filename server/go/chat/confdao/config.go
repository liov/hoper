package confdao

import (
	"github.com/hopeio/dora/initialize/gormdb"
	"github.com/hopeio/dora/initialize/log"
	"github.com/hopeio/dora/initialize/redis"
	"github.com/hopeio/dora/initialize/ristretto"
	"github.com/hopeio/dora/initialize/server"
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
