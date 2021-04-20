package conf

import (
	"runtime"
	"time"

	"github.com/liov/hoper/go/v2/tailmon/initialize"
	"github.com/spf13/viper"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &RedisConfig{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    initialize.ServerConfig
	Mail      initialize.MailConfig
	GORMDB    initialize.DatabaseConfig
	Redis     initialize.RedisConfig
	Cache     initialize.CacheConfig
	Log       initialize.LogConfig
	Viper     *viper.Viper
}

var Conf = &config{}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}
