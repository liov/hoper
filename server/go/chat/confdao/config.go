package confdao

import (
	"github.com/hopeio/initialize"
	"github.com/hopeio/initialize/conf_dao/server"
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
