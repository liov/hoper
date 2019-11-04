package initialize

import (
	"flag"
	"reflect"
	"strings"

	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
)

//var closes = []interface{}{log.Sync}

type Init struct {
	Env        string
	Module     string
	NoInit     []string
	Additional string //附加配置，不对外公开的的配置
	conf       needInit
	dao        dao
	//closes     []interface{}
}

type needInit interface {
	Custom()
}

type dao interface {
	Close()
	needInit
}

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf needInit, dao dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}
	init := &Init{conf: conf, dao: dao}
	init.config()
	value := reflect.ValueOf(init)
	noInit := strings.Join(init.NoInit, " ")
	typeOf := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		if strings.Contains(noInit, typeOf.Method(i).Name[2:]) {
			continue
		}
		if typeOf.Method(i).Type.NumOut() > 0 && dao == nil {
			continue
		}

		if res := value.Method(i).Call(nil); res != nil && len(res) > 0 {
			daoValue := reflect.ValueOf(dao).Elem()
			for j := range res {
				if res[j].IsValid() {
					reflect3.SetFieldValue(daoValue, res[j])
				}
			}
		}
	}
	if dao != nil {
		dao.Custom()
	}

	return func() {
		if dao != nil {
			dao.Close()
		}
		log.Sync()
		/*for _, f := range closes {
			res := reflect.ValueOf(f).Call(nil)
			if len(res) > 0 && res[0].IsValid() {
				log.Error(res[0].Interface())
			}
		}*/
	}
}
