package confdao

import (
	"github.com/hopeio/zeta/initialize/gormdb"
	"github.com/hopeio/zeta/initialize/log"
	"github.com/hopeio/zeta/initialize/redis"
	"github.com/hopeio/zeta/initialize/ristretto"
	"github.com/hopeio/zeta/initialize/server"
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
