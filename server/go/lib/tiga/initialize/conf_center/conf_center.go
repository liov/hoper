package conf_center

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/conf_center/local"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/conf_center/nacos"
	v2 "github.com/actliboy/hoper/server/go/lib/tiga/initialize/conf_center/nacos/v2"
)

const (
	Local   = "local"
	Nacos   = "nacos"
	Nacosv2 = "nacosv2"
)

type ConfigCenter interface {
	HandleConfig(func([]byte)) error
}

type ConfigCenterConfig struct {
	ConfigType string
	Watch      bool
	Nacos      *nacos.Nacos
	Nacosv2    *v2.Nacos
	Local      *local.Local
	/*	Etcd   *etcd.Etcd
		Apollo *apollo.Apollo*/
}

func (c *ConfigCenterConfig) ConfigCenter(model string, debug bool) ConfigCenter {
	if c.ConfigType == Nacos && c.Nacos != nil {
		c.Nacos.Watch = c.Watch
		return c.Nacos
	}

	if c.ConfigType == Nacosv2 && c.Nacosv2 != nil {
		return c.Nacosv2
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
