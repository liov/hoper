package global

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry"
	"github.com/hopeio/initialize/rootconf"
	"github.com/hopeio/utils/os/fs"
	timei "github.com/hopeio/utils/time"
	"time"
)

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &RedisConfig{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	PageSize      int8
	Volume        fs.Dir
	SiteURL       string
	QrCodeSaveDir fs.Dir //二维码保存路径
	FontSaveDir   fs.Dir //字体保存路径
	User          Config
	Moment        Moment
	Upload        Upload
	Server        cherry.Server
}

func (c *config) BeforeInject() {
	c.User.TokenMaxAge = timei.Day
}

func (c *config) AfterInject() {
	c.User.TokenMaxAge = timei.StdDuration(c.User.TokenMaxAge, time.Hour)
	c.User.TokenSecretBytes = []byte(c.User.TokenSecret)
}

func (c *config) AfterInjectWithRoot(rootconfig *rootconf.RootConfig) {
	if !rootconfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

type Config struct {
	PassSalt string
	// 天数
	TokenMaxAge      time.Duration
	TokenSecret      string
	TokenSecretBytes []byte

	LuosimaoSuperPW   string
	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	Limit UserLimit
}

type UserLimit struct {
	// 用户名的最大长度
	MaxUserNameLen uint8

	// 用户名的最小长度
	MinUserNameLen uint8

	// 密码的最大长度
	MaxPassLen uint8

	// 密码的最小长度
	MinPassLen uint8

	// 个性签名最大长度
	MaxSignatureLen uint8

	// 居住地的最大长度
	MaxLocationLen uint8

	//个人简介的最大长度
	MaxIntroduceLen uint8
}

type ContentLimit struct {
	SecondLimitKey, MinuteLimitKey, DayLimitKey       string
	SecondLimitCount, MinuteLimitCount, DayLimitCount int64
}

type Moment struct {
	MaxContentLen int
	Limit         ContentLimit
}

type Upload struct {
	Volume fs.Dir

	UploadDir      fs.Dir
	UploadMaxSize  int64
	UploadAllowExt []string
}
