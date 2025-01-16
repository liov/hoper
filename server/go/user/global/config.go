package global

import (
	"github.com/hopeio/cherry"
	timei "github.com/hopeio/utils/time"
	"time"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &RedisConfig{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    cherry.Server
}

var Conf = &config{}

func (c *config) BeforeInject() {
	c.Customize.TokenMaxAge = timei.Day
}

func (c *config) AfterInject() {
	c.Customize.TokenMaxAge = timei.StdDuration(c.Customize.TokenMaxAge, time.Hour)
	c.Customize.TokenSecretBytes = []byte(c.Customize.TokenSecret)
}
