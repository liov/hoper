package conf_center

import (
	"github.com/liov/hoper/v2/utils/configor/nacos"
)

type Nacos nacos.Config

// 从nacos拉取配置并返回nacos client
func (cc *Nacos) SetConfig(handle func([]byte)) error {
	nacosClient := (*nacos.Config)(cc).NewClient()
	err := nacosClient.GetConfigAllInfoHandle(handle)
	go nacosClient.Listener(handle)
	return err
}
