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
	Addr       string
	AppId      string `json:"appId"`
	Cluster    string `json:"cluster"`
	IP         string `json:"ip"`
	InitConfig apollo.SpecialConfig
	NameSpace  []string
}

func (conf *ApolloConfig) Generate() *apollo.Server {
	return apollo.New(conf.Addr, conf.AppId, conf.Cluster, conf.IP, conf.NameSpace, conf.InitConfig)
}

func (init *Init) P9Apollo() *apollo.Server {
	//配置更新时创建新的sever关闭旧的
	dao := reflect.ValueOf(init.dao)
	for i := 0; i < dao.NumField(); i++ {
		if dao.Field(i).Type() == reflect.TypeOf(&apollo.Server{}) {
			s := dao.Field(i).Interface()
			if s != nil {
				s.(*apollo.Server).Close()
			}

		}
	}

	conf := &ApolloConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return nil
	}
	conf.NameSpace = append(conf.NameSpace, "initialize")
	cCopy := init.conf
	conf.InitConfig = apollo.SpecialConfig{NameSpace: "initialize", Callback: func(m map[string]string) {
		apolloConfigEnable(cCopy, m)
	}}
	return conf.Generate()
}

//优先级应该最低，要更新配置
func (init *Init) apollo() {
	conf := &ApolloConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return
	}
	c := apollo.NewConfig(conf.Addr, conf.AppId, conf.Cluster, conf.IP)
	aConf, err := c.GetInitConfig("initialize")
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
