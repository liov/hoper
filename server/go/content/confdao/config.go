package confdao

import (
	"github.com/hopeio/tiga/initialize/conf_dao/log"
	"github.com/hopeio/tiga/initialize/conf_dao/server"
	"github.com/liov/hoper/server/go/content/model"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &Config{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.Config
	Log       log.Config
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

func (c *config) InitBeforeInject() {

}

func (c *config) InitAfterInject() {

}
