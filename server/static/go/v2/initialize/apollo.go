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
	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP, conf.NameSpace)
}

//优先级应该最低，要更新配置
func (init *Init) P9Apollo() *apollo.Server {
	conf := &ApolloConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return nil
	}
	return conf.Generate()
}
