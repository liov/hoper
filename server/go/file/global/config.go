package global

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/gox/os/fs"
)

type config struct {
	//自定义的配置
	Customize Config
	Server    cherry.Server
}

func (c *config) BeforeInject() {
	c.Customize.UploadMaxSize = 1024 * 1024 * 1024
}

func (c *config) AfterInject() {
	c.Customize.UploadMaxSize = c.Customize.UploadMaxSize * 1024 * 1024
}

type Config struct {
	Volume fs.Dir

	UploadDir      fs.Dir
	UploadMaxSize  int64
	UploadAllowExt []string
}
