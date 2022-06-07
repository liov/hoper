package conf_center

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/conf_center/local"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/conf_center/nacos"
	local2 "github.com/actliboy/hoper/server/go/lib/utils/configor/local"
)

type ConfigCenter interface {
	HandleConfig(func([]byte)) error
}

type ConfigCenterConfig struct {
	Watch bool
	Nacos *nacos.Nacos
	Local *local.Local
	/*	Etcd   *etcd.Etcd
		Apollo *apollo.Apollo*/
}

type ConfigCenterEnvConfig struct {
	NacosTenant string
	//本地配置，特定文件名,启用文件搜寻查找
	LocalConfigName string
	EtcdKey         string
}

func (c *ConfigCenterConfig) ConfigCenter(ccec ConfigCenterEnvConfig, model string, debug bool) ConfigCenter {
	if c.Nacos != nil && ccec.NacosTenant != "" {
		c.Nacos.DataId = model
		c.Nacos.Tenant = ccec.NacosTenant
		c.Nacos.Watch = c.Watch
		return c.Nacos
	}
	/*	if c.Etcd != nil && ccec.EtcdKey != "" {
		c.Etcd.Key = ccec.EtcdKey
		c.Etcd.Watch = c.Watch
		return c.Etcd
	}*/
	if ccec.LocalConfigName != "" {
		if c.Local != nil {
			c.Local.LocalConfigName = ccec.LocalConfigName
			c.Local.AutoReload = c.Watch
			return c.Local
		} else {
			return &local.Local{Config: local2.Config{
				Debug:      debug,
				AutoReload: c.Watch,
			},
				LocalConfigName: ccec.LocalConfigName}
		}
	}
	panic("没有设置配置中心")
}
