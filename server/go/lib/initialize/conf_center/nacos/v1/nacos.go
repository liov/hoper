package v1

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type Nacos struct {
	vo.NacosClientParam
	vo.ConfigParam
}

// 从nacos拉取配置并返回nacos client
func (cc *Nacos) HandleConfig(handle func([]byte)) error {
	client, err := clients.NewConfigClient(cc.NacosClientParam)
	if err != nil {
		log.Fatal(err)
	}
	cc.OnChange = func(namespace, group, dataId, data string) {
		handle([]byte(data))
	}
	err = client.ListenConfig(cc.ConfigParam)
	return err
}
