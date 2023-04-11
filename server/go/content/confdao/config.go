package confdao

import (
	"github.com/hopeio/pandora/initialize/gormdb"
	"github.com/hopeio/pandora/initialize/log"
	"github.com/hopeio/pandora/initialize/mail"
	"github.com/hopeio/pandora/initialize/redis"
	"github.com/hopeio/pandora/initialize/ristretto"
	"github.com/hopeio/pandora/initialize/server"
	"runtime"
	"time"

	"github.com/liov/hoper/server/go/mod/content/model"
	"github.com/spf13/viper"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &Config{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.ServerConfig
	Mail      mail.MailConfig
	GORMDB    gormdb.DatabaseConfig
	Redis     redis.Config
	Cache     ristretto.CacheConfig
	Log       log.LogConfig
	Viper     *viper.Viper
}

var Conf = &config{
	Customize: serverConfig{
		Moment: Moment{
			Limit: Limit{
				SecondLimitKey: model.MomentSecondLimitKey,
				MinuteLimitKey: model.MomentMinuteLimitKey,
				DayLimitKey:    model.MomentDayLimitKey,
			},
		},
	},
}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
		c.Customize.LuosimaoAPIKey = ""
	}

	c.Server.ReadTimeout = c.Server.ReadTimeout * time.Second
	c.Server.WriteTimeout = c.Server.WriteTimeout * time.Second
}
