package initialize

import (
	"github.com/liov/hoper/go/v2/tools/apollo"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

type ApolloConfig struct {
	Addr      string
	AppId     string `json:"appId"`
	Cluster   string `json:"cluster"`
	IP        string `json:"ip"`
	NameSpace []string
}

func (conf *ApolloConfig) Generate() *apollo.Server {
	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP)
}

func (i *Init) P0Apollo() *apollo.Server {
	conf := &ApolloConfig{}
	if exist := reflect3.GetFieldValue(i.conf, conf); !exist {
		return nil
	}
	return conf.Generate()
}
