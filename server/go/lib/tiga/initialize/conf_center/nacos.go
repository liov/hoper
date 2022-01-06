package conf_center

import (
	"github.com/actliboy/hoper/server/go/lib/utils/configor/nacos"
)

type Nacos struct {
	nacos.Config
}

// 从nacos拉取配置并返回nacos client
func (cc *Nacos) HandleConfig(handle func([]byte)) error {
	nacosClient := cc.NewClient()
	err := nacosClient.GetConfigAllInfoHandle(handle)
	go nacosClient.Listener(handle)
	return err
}
