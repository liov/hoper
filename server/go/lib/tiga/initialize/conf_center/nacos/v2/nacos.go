package v2

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/v2/common/file"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
)

type Nacos struct {
	vo.NacosClientParam
	vo.ConfigParam
}

// 由于库暴露的api可diy性太差，废弃
func (cc *Nacos) HandleConfig(handle func([]byte)) error {
	client, err := clients.NewConfigClient(cc.NacosClientParam)
	if err != nil {
		log.Fatal(err)
	}
	config, err := client.GetConfig(cc.ConfigParam)
	if err != nil {
		log.Fatal(err)
	}
	// nacos-go-sdk的问题，首次拉取的配置缓存在cache目录，listen拉取的缓存在cache/config，listen是异步的，如果要先同步获取配置且不在未更改配置的情况下触发listen的Onchange，就要把配置写进listen的目录，来回读取写入，浪费性能
	cacheDir := file.GetCurrentPath() + string(os.PathSeparator) + "cache/config"
	cacheKey := util.GetConfigCacheKey(cc.DataId, cc.Group, cc.ClientConfig.NamespaceId)
	cache.WriteConfigToFile(cacheKey, cacheDir, config)
	handle([]byte(config))
	cc.OnChange = func(namespace, group, dataId, data string) {
		handle([]byte(data))
	}

	return client.ListenConfig(cc.ConfigParam)
}
