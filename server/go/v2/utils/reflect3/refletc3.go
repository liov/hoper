package reflect3

import (
	"fmt"
	"reflect"

	"github.com/modern-go/reflect2"
)

func TypeInfo(v interface{}) {
	value := reflect.ValueOf(v).Elem()
	type2 := reflect2.TypeOf(&value)
	rtype := type2.(reflect2.PtrType).Elem().(reflect2.StructType)
	typField := rtype.FieldByName("typ")
	typTyp := typField.Type().(reflect2.PtrType).Elem().(reflect2.StructType)
	typV := typField.Get(&value)
	typV = reflect.ValueOf(typV).Elem().Interface()
	for i := 0; i < typTyp.NumField(); i++ {
		field := typTyp.Field(i)
		v := field.Get(typV)
		v = reflect.ValueOf(v).Elem().Interface()
		fmt.Printf("字段名：%v,字段值：%v\n", field.Name(), v)
	}
}

func Set(o interface{}, field string, v interface{}) {
	t := reflect2.TypeOf(o)
	if field == "" {
		t.Set(o, v)
	} else {
		f := t.(reflect2.PtrType).Elem().(reflect2.StructType).FieldByName(field)
		f.Set(o, v)
	}
}

func OriginalType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return OriginalType(typ.Elem())
	}
	return typ
}
