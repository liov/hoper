package conf

import (
	"runtime"
	"time"

	"github.com/liov/hoper/go/v2/content/model"
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

var Conf = &config{
	Customize: serverConfig{
		Moment: Moment{
			Limit: Limit{
				SecondLimit: model.MomentSecondLimitKey,
				MinuteLimit: model.MomentMinuteLimitKey,
				DayLimit:    model.MomentDayLimitKey,
			},
		},
	},
}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
		c.Customize.LuosimaoAPIKey = ""
	}

	c.Server.ReadTimeout = c.Server.ReadTimeout * time.Second
	c.Server.WriteTimeout = c.Server.WriteTimeout * time.Second
}
