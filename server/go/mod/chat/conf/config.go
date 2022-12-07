package conf

import (
	"github.com/liov/hoper/server/go/lib/initialize/db"
	"github.com/liov/hoper/server/go/lib/initialize/log"
	"github.com/liov/hoper/server/go/lib/initialize/redis"
	"github.com/liov/hoper/server/go/lib/initialize/server"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.ServerConfig
	GORMDB    db.DatabaseConfig
	Redis     redis.Config
	Cache     cache_ristretto.CacheConfig
	Log       log.LogConfig
}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

}

type serverConfig struct {
	Host string
}
