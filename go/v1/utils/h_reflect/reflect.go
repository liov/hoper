package h_reflect

import (
	"reflect"
)

func ContainType()  {

}

func GetExpectTypeValue(src interface{}, dst interface{}) bool {
	srcValue:= reflect.ValueOf(src).Elem()
	for i := 0; i < srcValue.NumField(); i++ {
		if srcValue.Field(i).Type() == reflect.ValueOf(dst).Elem().Type(){
			reflect.ValueOf(dst).Elem().Set(srcValue.Field(i))
			return true
		}
	}
	return false
}