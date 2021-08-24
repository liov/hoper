package conf

import (
	"github.com/liov/hoper/server/go/lib/tiga/initialize"
	"github.com/liov/hoper/server/go/lib/tiga/initialize/inject_dao"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    initialize.ServerConfig
	GORMDB    inject_dao.DatabaseConfig
	Redis     inject_dao.RedisConfig
	Cache     inject_dao.CacheConfig
	Log       initialize.LogConfig
}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

}

type serverConfig struct {
	Host string
}
