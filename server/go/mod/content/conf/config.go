package conf

import (
	"github.com/actliboy/hoper/server/go/lib/initialize/cache/ristretto"
	"github.com/actliboy/hoper/server/go/lib/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/initialize/log"
	"github.com/actliboy/hoper/server/go/lib/initialize/mail"
	"github.com/actliboy/hoper/server/go/lib/initialize/redis"
	"github.com/actliboy/hoper/server/go/lib/initialize/server"
	"runtime"
	"time"

	"github.com/actliboy/hoper/server/go/mod/content/model"
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
	GORMDB    db.DatabaseConfig
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
