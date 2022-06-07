package v1

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
	service, _ := c.GetService("user")
	log.Println(service.Name)
}

func TestBytes(t *testing.T) {
	var data []byte
	log.Println(string(data))
}
