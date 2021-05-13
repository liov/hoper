package initialize

import (
	"github.com/liov/hoper/v2/utils/configor/apollo"
	"github.com/liov/hoper/v2/utils/log"
	"github.com/pelletier/go-toml"
	"reflect"
	"strings"
)

func (init *Init) getApolloClient() *apollo.Client {
	return nil
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
