package initialize

import (
	"github.com/liov/hoper/go/v2/utils/configor/nacos"
	"github.com/liov/hoper/go/v2/utils/log"
)

func (init *Init) Register() {
	if init.ConfigCenter == nil {
		return
	}
	svcName := init.BasicConfig.Module
	_, err := init.ConfigCenter.GetService(svcName)
	serviceConfig := init.GetServiceConfig()
	if err != nil {
		err = init.ConfigCenter.CreateService(svcName, &nacos.Metadata{
			Domain: serviceConfig.Domain,
			Port:   serviceConfig.Port,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	err = init.ConfigCenter.RegisterInstance(serviceConfig.Port, svcName)
	if err != nil {
		log.Fatal(err)
	}
}
