package config

import (
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	"github.com/liov/hoper/go/v2/utils/fs"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server   initialize.ServerConfig
	Database initialize.DatabaseConfig
	Redis    initialize.RedisConfig
	Log      initialize.LogConfig
}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.UploadMaxSize = c.Customize.UploadMaxSize * 1024 * 1024
}

type serverConfig struct {
	Volume fs.Dir


	UploadDir      fs.Dir
	UploadMaxSize  int64
	UploadAllowExt []string

	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSaveDir fs.Dir //二维码保存路径
	PrefixUrl     string
	FontSaveDir   fs.Dir //字体保存路径
}