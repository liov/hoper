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
			DataId: "test",
			Group:  "DEFAULT_GROUP",
		},
		MD5:   "",
		close: nil,
	}
	config, _ := c.GetService("user")
	log.Println(string(config))
}
