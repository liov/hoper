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
	Server   initialize.ServerConfig
	Mail     initialize.MailConfig
	Database initialize.DatabaseConfig
	Redis    initialize.RedisConfig
	Log      initialize.LogConfig
	Viper     *viper.Viper
}

var Conf = &config{}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.UploadMaxSize = c.Customize.UploadMaxSize * 1024 * 1024
	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}
