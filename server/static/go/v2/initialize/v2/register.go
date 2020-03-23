package initialize

import (
	"github.com/liov/hoper/go/v2/tools/nacos"
	"github.com/liov/hoper/go/v2/utils/log"
)

func (init *Init) Register() {
	svcName := init.NacosConfig.DataId
	_, err := init.NacosConfig.GetService(svcName)
	if err != nil {
		err = init.NacosConfig.CreateService(svcName, &nacos.Metadata{Domain: init.GetServiceDomain()})
		if err != nil {
			log.Fatal(err)
		}
	}
	err = init.NacosConfig.RegisterInstance(init.GetServicePort(), svcName)
	if err != nil {
		log.Fatal(err)
	}
}
