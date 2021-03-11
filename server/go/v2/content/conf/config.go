package conf

import (
	"runtime"
	"time"

	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/initialize"
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
	Database  initialize.DatabaseConfig
	Redis     initialize.RedisConfig
	Log       initialize.LogConfig
	Consul    initialize.EtcdConfig
}

var Conf = &config{
	Customize: serverConfig{
		Moment: Moment{
			Limit:Limit{
				SecondLimit:      model.MomentSecondLimitKey,
				MinuteLimit:      model.MomentMinuteLimitKey,
				DayLimit:         model.MomentDayLimitKey,
			},
		},
	},
}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
		c.Customize.LuosimaoAPIKey = ""
	}

	c.Customize.UploadMaxSize = c.Customize.UploadMaxSize * 1024 * 1024
	c.Server.ReadTimeout = c.Server.ReadTimeout * time.Second
	c.Server.WriteTimeout = c.Server.WriteTimeout * time.Second
}
