package conf

import (
	"github.com/liov/hoper/v2/tiga/initialize"
	"github.com/liov/hoper/v2/tiga/initialize/inject_dao"
	"github.com/liov/hoper/v2/utils/fs"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    initialize.ServerConfig
	GORMDB    inject_dao.DatabaseConfig
	Redis     inject_dao.RedisConfig
	Cache     inject_dao.CacheConfig
	Log       initialize.LogConfig
}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.UploadMaxSize = c.Customize.UploadMaxSize * 1024 * 1024
}

type serverConfig struct {
	Volume fs.Dir

	UploadUrlPrefix string
	UploadDir       fs.Dir
	UploadMaxSize   int64
	UploadAllowExt  []string

	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSaveDir fs.Dir //二维码保存路径
	FontSaveDir   fs.Dir //字体保存路径
}
