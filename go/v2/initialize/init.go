package initialize

import (
	"flag"
	"reflect"
	"strings"

	"github.com/liov/hoper/go/v2/utils/reflect3"
)

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
)

type Init struct {
	Env    string
	Module string
	NoInit []string
	conf   needInit
	dao    needInit
}

type needInit interface {
	Custom()
}

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf needInit, dao needInit) {
	if !flag.Parsed(){
		flag.Parse()
	}
	init := &Init{conf: conf,dao:dao}
	init.config()
	value := reflect.ValueOf(init)
	noInit := strings.Join(init.NoInit, " ")
	typeOf := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		if strings.Contains(noInit, typeOf.Method(i).Name[2:]) {
			continue
		}
		if res:=value.Method(i).Call(nil);res!=nil && len(res)>0{
			daoValue:= reflect.ValueOf(dao).Elem()
			for j := range res{
				if res[j].IsValid(){
					reflect3.SetFieldValue(daoValue,res[j])
				}
			}
		}
	}
	if dao != nil {
		dao.Custom()
	}
}
