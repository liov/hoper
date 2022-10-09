package conf_center

import (
	"github.com/actliboy/hoper/server/go/lib/initialize/conf_center/http"
	"github.com/actliboy/hoper/server/go/lib/initialize/conf_center/local"
	"github.com/actliboy/hoper/server/go/lib/initialize/conf_center/nacos/v2"
)

const (
	Local = "local"
	Nacos = "nacos"
	Http  = "http"
)

type ConfigCenter interface {
	HandleConfig(func([]byte)) error
}

type ConfigCenterConfig struct {
	ConfigType string
	Watch      bool
	Nacos      *v2.Nacos
	Local      *local.Local
	Http       *http.Config
	/*	Etcd   *etcd.Etcd
		Apollo *apollo.Apollo*/
}

func (c *ConfigCenterConfig) ConfigCenter(model string, debug bool) ConfigCenter {
	if c.ConfigType == Http && c.Http != nil {
		return c.Http
	}

	if c.ConfigType == Nacos && c.Nacos != nil {
		return c.Nacos
	}
	/*	if c.Etcd != nil && ccec.EtcdKey != "" {
		c.Etcd.Key = ccec.EtcdKey
		c.Etcd.Watch = c.Watch
		return c.Etcd
	}*/
	if c.ConfigType == Local && c.Local != nil {
		c.Local.Debug = debug
		c.Local.AutoReload = c.Watch
		return c.Local
	}
	return c.Local
}
