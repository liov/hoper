package conf

import (
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    initialize.ServerConfig
	GORMDB    initialize.DatabaseConfig
	Redis     initialize.RedisConfig
	Cache     initialize.CacheConfig
	Log       initialize.LogConfig
}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
	}

}

type serverConfig struct {
	Host string
}
