package global

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/utils/os/fs"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    cherry.Server
}

func (c *config) BeforeInject() {
	c.Customize.UploadMaxSize = 1024 * 1024 * 1024
}

func (c *config) AfterInject() {
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
	FontSaveDir   fs.Dir //字体保存路径
}
