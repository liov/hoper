package confdao

import (
	"github.com/hopeio/tiga/initialize/conf_dao/log"
	"github.com/hopeio/tiga/initialize/conf_dao/server"
	"github.com/hopeio/tiga/initialize/conf_dao/viper"
	timei "github.com/hopeio/tiga/utils/time"
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
	Log       log.Config
	Viper     viper.Config
}

var Conf = &config{}

func (c *config) InitBeforeInject() {
	c.Customize.TokenMaxAge = timei.Day
}

func (c *config) InitAfterInject() {
	c.Customize.TokenMaxAge = timei.StdDuration(c.Customize.TokenMaxAge, time.Hour)
	c.Customize.TokenSecretBytes = []byte(c.Customize.TokenSecret)
}
