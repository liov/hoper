package initialize

import (
	"reflect"
	"strings"
)

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
)

type Init struct {
	Env    string
	NoInit []string
}

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf Config, customInit func()) {
	init := &Init{}
	init.config(conf)
	value := reflect.ValueOf(init)
	noInit := strings.Join(init.NoInit, " ")
	typeOf := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		if strings.Contains(noInit, typeOf.Method(i).Name[2:]) {
			continue
		}
		value.Method(i).Call([]reflect.Value{reflect.ValueOf(conf)})
	}
	if customInit != nil {
		customInit()
	}
}
