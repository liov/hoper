package conf

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/cache"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/log"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/redis"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/server"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.ServerConfig
	GORMDB    db.DatabaseConfig
	Redis     redis.RedisConfig
	Cache     cache.CacheConfig
	Log       log.LogConfig
}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

}

type serverConfig struct {
	Host string
}
