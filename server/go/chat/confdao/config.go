package confdao

import (
	"github.com/hopeio/cherry/initialize"
	"github.com/hopeio/cherry/initialize/conf_dao/server"
)

var Conf = &config{}

type config struct {
	initialize.EmbeddedPresets
	//自定义的配置
	Customize serverConfig
	Server    server.Config
}

type serverConfig struct {
	Host string
}
