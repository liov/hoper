package confdao

import (
	"github.com/hopeio/cherry/initialize/conf_dao/log"
	"github.com/hopeio/cherry/initialize/conf_dao/server"
	"github.com/hopeio/cherry/utils/io/fs"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.Config
	Log       log.Config
}

func (c *config) InitBeforeInject() {
	c.Customize.UploadMaxSize = 1024 * 1024 * 1024
}

func (c *config) InitAfterInject() {
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
