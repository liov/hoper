package reflect3

import (
	"reflect"
)

func ContainType()  {

}
//获取子类型的值
//参数父类型，子类型的指针
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