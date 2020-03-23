package nacos

import (
	"log"
	"testing"
)

func TestNacos(t *testing.T) {
	c := Client{
		Config: &Config{
			Addr:   "192.168.1.212:9001",
			Tenant: "",
			DataId: "user",
			Group:  "DEFAULT_GROUP",
		},
		MD5:   "",
		close: nil,
	}
	service, _ := c.GetService()
	log.Println(service.Name)
}
