package reflecti

import (
	"fmt"
	"reflect"

	"github.com/liov/hoper/go/v2/utils/log"
)

func ContainType() {

}

//获取子类型的值
//参数父类型，子类型的指针
func GetFieldValue(src interface{}, dst interface{}) bool {
	srcValue := reflect.ValueOf(src).Elem()
	for i := 0; i < srcValue.NumField(); i++ {
		if srcValue.Field(i).Type() == reflect.ValueOf(dst).Elem().Type() {
			reflect.ValueOf(dst).Elem().Set(srcValue.Field(i))
			return true
		}
	}
	return false
}

func SetField(src interface{}, sub interface{}) bool {
	srcValue := reflect.ValueOf(src).Elem()
	subValue := reflect.ValueOf(sub)
	SetFieldValue(srcValue, subValue)
	subValue = subValue.Elem()
	return SetFieldValue(srcValue, subValue)
}

func SetFieldValue(srcValue reflect.Value, subValue reflect.Value) bool {
	for i := 0; i < srcValue.NumField(); i++ {
		if srcValue.Field(i).Type() == subValue.Type() {
			srcValue.Field(i).Set(subValue)
			return true
		}
	}
	return false
}

func setField(obj interface{}, name string, value interface{}) error {
	structData := reflect.ValueOf(obj).Elem()
	fieldValue := structData.FieldByName(name)

	if !fieldValue.IsValid() {
		return fmt.Errorf("utils.setField() No such field: %s in obj ", name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value ", name)
	}

	fieldType := fieldValue.Type()
	val := reflect.ValueOf(value)

	valTypeStr := val.Type().String()
	fieldTypeStr := fieldType.String()
	if valTypeStr == "float64" && fieldTypeStr == "int" {
		val = val.Convert(fieldType)
	} else if fieldType != val.Type() {
		return fmt.Errorf("Provided value type " + valTypeStr + " didn't match obj field type " + fieldTypeStr)
	}
	fieldValue.Set(val)
	return nil
}

// SetStructByJSON 由json对象生成 struct
func SetStructByJSON(obj interface{}, mapData map[string]interface{}) error {
	for key, value := range mapData {
		if err := setField(obj, key, value); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
