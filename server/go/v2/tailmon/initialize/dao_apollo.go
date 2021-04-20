package initialize

import (
	"reflect"
	"strings"

	"github.com/liov/hoper/go/v2/utils/configor/apollo"
	"github.com/liov/hoper/go/v2/utils/log"
	reflecti "github.com/liov/hoper/go/v2/utils/reflect"
	"github.com/pelletier/go-toml"
)

func (conf *ApolloConfig) generate() *apollo.Client {
	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP, conf.NameSpace, conf.InitConfig)
}

func (conf *ApolloConfig) Generate() interface{} {
	return conf.generate()
}

func (init *Init) P0Apollo() *apollo.Client {
	conf := &ApolloConfig{}
	if exist := reflecti.GetFieldValue(init.conf, conf); !exist {
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
	conf.NameSpace = append(conf.NameSpace, conf.InitNameSpace)

	return conf.generate()
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
