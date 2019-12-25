package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/liov/hoper/go/v2/tools/apollo"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/pelletier/go-toml"
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

func (init *Init) P9Apollo() *apollo.Server {
	conf := &ApolloConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return nil
	}
	return conf.Generate()
}

//优先级应该最低，要更新配置
func (init *Init) apollo() {
	s := init.P9Apollo()
	aConf, err := s.GetInitConfig("initialize")
	if err != nil {
		panic(err)
	}
	apolloConfigEnable(init.conf, aConf)
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
	fmt.Println(conf)
}
