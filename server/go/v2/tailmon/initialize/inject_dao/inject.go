package initialize

import (
	"reflect"

	reflecti "github.com/liov/hoper/go/v2/utils/reflect"
)

type Inject struct{
	Init reflect.Value
}

func (in *Inject) SetConf(conf interface{}) bool{
	return reflecti.SetFieldValue(in.Init, reflect.ValueOf(conf))
}