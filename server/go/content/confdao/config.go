package confdao

import (
	"github.com/hopeio/zeta/initialize/gormdb"
	"github.com/hopeio/zeta/initialize/log"
	"github.com/hopeio/zeta/initialize/mail"
	"github.com/hopeio/zeta/initialize/redis"
	"github.com/hopeio/zeta/initialize/ristretto"
	"github.com/hopeio/zeta/initialize/server"
	"runtime"
	"time"

	"github.com/actliboy/hoper/server/go/content/model"
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
