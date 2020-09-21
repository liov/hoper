package initialize

import (
	"reflect"
	"strings"

	"github.com/liov/hoper/go/v2/utils/apollo"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/pelletier/go-toml"
)

type ApolloConfig struct {
	Addr       string
	AppId      string `json:"appId"`
	Cluster    string `json:"cluster"`
	IP         string `json:"ip"`
	InitConfig apollo.SpecialConfig
	NameSpace  []string
}

func (conf *ApolloConfig) Generate() *apollo.Client {
	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP, conf.NameSpace, conf.InitConfig)
}

func (init *Init) P0Apollo() *apollo.Client {
	conf := &ApolloConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return nil
	}
	//初始化更新配置，这里不需要，开启实时更新时初始化会更新一次
	/*	c := apollo.NewConfig(conf.Addr, conf.AppId, conf.Cluster, conf.IP)
		aConf, err := c.GetInitConfig(InitKey)
		if err != nil {
			panic(err)
		}
		apolloConfigEnable(init.conf, aConf)*/
	//监听指定namespace的更新
	conf.NameSpace = append(conf.NameSpace, InitKey)
	cCopy := init.conf
	dCopy := init.dao
	conf.InitConfig = apollo.SpecialConfig{NameSpace: InitKey, Callback: func(m map[string]string) {
		apolloConfigEnable(cCopy, m)
		Refresh(cCopy, dCopy)
	}}
	return conf.Generate()
}

func apolloConfigEnable(conf interface{}, aConf map[string]string) {
	confValue := reflect.ValueOf(conf).Elem()
	for k, v := range aConf {
		field := confValue.FieldByNameFunc(func(s string) bool {
			return strings.ToUpper(s) == strings.ToUpper(k)
		})
		subConf := field.Addr().Interface()
		err := toml.Unmarshal([]byte(v), subConf)
		//err := json.Unmarshal([]byte(v),subConf)
		if err != nil {
			log.Error(err)
		}
	}
	log.Debug(conf)
}
