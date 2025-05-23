package global

import (
	"github.com/hopeio/cherry"
	"github.com/liov/hoper/server/go/content/model"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &Config{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize Config
	Server    cherry.Server
}

func (c *config) BeforeInject() {
	c.Customize = Config{
		Moment: Moment{
			Limit: Limit{
				SecondLimitKey: model.MomentSecondLimitKey,
				MinuteLimitKey: model.MomentMinuteLimitKey,
				DayLimitKey:    model.MomentDayLimitKey,
			},
		},
	}
}

func (c *config) AfterInject() {

}
