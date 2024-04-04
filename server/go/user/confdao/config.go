package confdao

import (
	"github.com/hopeio/tiga/initialize/conf_dao/log"
	"github.com/hopeio/tiga/initialize/conf_dao/server"
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
	Server    server.Config
	Log       log.Config
	Viper     *viper.Viper
}

var Conf = &config{}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}
