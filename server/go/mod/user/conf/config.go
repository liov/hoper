package conf

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/log"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/mail"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/redis"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/server"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &RedisConfig{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.ServerConfig
	Mail      mail.MailConfig
	GORMDB    db.DatabaseConfig
	Redis     redis.RedisConfig
	Cache     cache_ristretto.CacheConfig
	Log       log.LogConfig
	Viper     *viper.Viper
}

var Conf = &config{}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}
