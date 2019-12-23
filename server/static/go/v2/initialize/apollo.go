package initialize

import (
	"github.com/liov/hoper/go/v2/tools/apollo"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P0Apollo() *apollo.Server {
	conf := ApolloConfig{}
	if exist := reflect3.GetFieldValue(i.conf, &conf); !exist {
		return nil
	}
	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP)
}
