package initialize

import (
	"github.com/liov/hoper/go/v2/tools/nacos"
)

func (init *Init) Register(client *nacos.Client) {
	svcName := init.NacosConfig.DataId
	_, err := client.GetService(svcName)
	if err != nil {
		client.CreateService(svcName)
	}
	client.RegisterInstance(init.GetServicePort(), svcName)
}
