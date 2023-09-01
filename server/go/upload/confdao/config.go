package confdao

import (
	"github.com/hopeio/lemon/initialize/basic_conf/log"
	"github.com/hopeio/lemon/initialize/basic_conf/server"
	"github.com/hopeio/lemon/utils/io/fs"
	"runtime"
)

var Conf = &config{}

type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.ServerConfig
	Log       log.LogConfig
}

func (c *config) Init() {
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
	FontSaveDir   fs.Dir //字体保存路径
}
