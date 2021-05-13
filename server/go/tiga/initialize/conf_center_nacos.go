package initialize

import (
	"github.com/liov/hoper/v2/utils/configor/nacos"
	"github.com/liov/hoper/v2/utils/log"
)

// 从nacos拉取配置并返回nacos client
func (init *Init) getNacosClient() *nacos.Client {
	nacosClient := init.ConfigCenter.NewClient()
	err := nacosClient.GetConfigAllInfoHandle(init.UnmarshalAndSet)
	if err != nil {
		log.Fatal(err)
	}
	return nacosClient
}
