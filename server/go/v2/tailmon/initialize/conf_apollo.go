package initialize

import (
	"github.com/liov/hoper/go/v2/utils/configor/apollo"
)

type ApolloConfig struct {
	Addr       string
	AppId      string `json:"appId"`
	Cluster    string `json:"cluster"`
	IP         string `json:"ip"`
	InitConfig apollo.SpecialConfig
	NameSpace  []string
	InitNameSpace string
}


func (conf *ApolloConfig) getApolloClient() *apollo.Client {
	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP, conf.NameSpace, conf.InitConfig)
}