package confdao

import (
	"github.com/hopeio/cherry/initialize/conf_dao/server"
	timei "github.com/hopeio/cherry/utils/time"
	"time"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &RedisConfig{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize serverConfig
	SendMail  SendMailConfig
	Server    server.Config
}

var Conf = &config{}

func (c *config) InitBeforeInject() {
	c.Customize.TokenMaxAge = timei.Day
}

func (c *config) InitAfterInject() {
	c.Customize.TokenMaxAge = timei.StdDuration(c.Customize.TokenMaxAge, time.Hour)
	c.Customize.TokenSecretBytes = []byte(c.Customize.TokenSecret)
}
