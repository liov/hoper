package global

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/initialize"
)

type config struct {
	initialize.EmbeddedPresets
	//自定义的配置
	Customize serverConfig
	Server    cherry.Server
}

type serverConfig struct {
	Host string
}
