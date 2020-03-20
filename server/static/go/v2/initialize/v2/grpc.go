package initialize

import (
	"github.com/liov/hoper/go/v2/tools/nacos"
)

func (init *Init) Register(client *nacos.Client) {
	client.GetService(init.GetServiceName())
}
