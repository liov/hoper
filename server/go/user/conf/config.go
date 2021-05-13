package conf

import (
	"github.com/liov/hoper/v2/tiga/initialize"
	"github.com/liov/hoper/v2/tiga/initialize/inject_dao"
	"runtime"
	"time"

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
	Mail      inject_dao.MailConfig
	GORMDB    inject_dao.DatabaseConfig
	Redis     inject_dao.RedisConfig
	Cache     inject_dao.CacheConfig
	Log       initialize.LogConfig
	Viper     *viper.Viper
}

var Conf = &config{}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}
