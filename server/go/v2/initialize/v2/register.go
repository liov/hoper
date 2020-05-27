package initialize

import (
	"github.com/liov/hoper/go/v2/tools/nacos"
	"github.com/liov/hoper/go/v2/utils/log"
)

func (init *Init) Register() {
	svcName := init.NacosConfig.DataId
	_, err := init.NacosConfig.GetService(svcName)
	serviceConfig := init.GetServiceConfig()
	if err != nil {
		err = init.NacosConfig.CreateService(svcName, &nacos.Metadata{
			Domain: serviceConfig.Domain,
			Port:   serviceConfig.Port,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	err = init.NacosConfig.RegisterInstance(serviceConfig.Port, svcName)
	if err != nil {
		log.Fatal(err)
	}
}
