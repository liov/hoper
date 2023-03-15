package apollo

import (
	"github.com/liov/hoper/server/go/lib/utils/configor/apollo"
)

type ApolloConfig struct {
	Addr          string
	AppId         string `json:"appId"`
	Cluster       string `json:"cluster"`
	IP            string `json:"ip"`
	InitConfig    apollo.SpecialConfig
	NameSpace     []string
	InitNameSpace string
}

func (conf *ApolloConfig) Build() *apollo.Client {
	//初始化更新配置，这里不需要，开启实时更新时初始化会更新一次
	/*	c := apollo.NewConfig(conf.Addr, conf.AppId, conf.Cluster, conf.IP)
		aConf, err := c.GetInitConfig(InitKey)
		if err != nil {
			panic(err)
		}
		apolloConfigEnable(init.conf, aConf)*/
	//监听指定namespace的更新
	conf.NameSpace = append(conf.NameSpace, conf.InitNameSpace)

	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP, conf.NameSpace, conf.InitConfig)
}

type ApolloClient struct {
	*apollo.Client
	Conf ApolloConfig
}

func (a *ApolloClient) Config() any {
	return &a.Conf
}

func (a *ApolloClient) SetEntity() {
	a.Client = a.Conf.Build()
}

func (c *ApolloClient) Close() error {
	return c.Client.Close()
}
