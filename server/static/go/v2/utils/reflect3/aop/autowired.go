package aop

import (
	"math/rand"
	"reflect"
	"strconv"

	"github.com/liov/hoper/go/v2/utils/mock"
)

func Autowired(v interface{}) {
	value := reflect.ValueOf(v)
	typMap := make(map[reflect.Type]struct{})
	autowired(value, typMap)
}

func autowired(value reflect.Value, typMap map[reflect.Type]struct{}) {
	typ := value.Type()
	typMap[typ] = struct{}{}
	typ = typ.Elem()
	value = value.Elem()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		tag := sf.Tag.Get("autowired")
		fieldValue := value.Field(i)
		switch fieldValue.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var r uint64
			if tag == "" {
				r = rand.Uint64()
			} else {
				r, _ = strconv.ParseUint(tag, 10, 64)
			}
			fieldValue.SetUint(r)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var r int64
			if tag == "" {
				r = rand.Int63()
			} else {
				r, _ = strconv.ParseInt(tag, 10, 64)
			}
			fieldValue.SetInt(r)
		case reflect.Float32, reflect.Float64:
			var r float64
			if tag == "" {
				r = rand.Float64()
			} else {
				r, _ = strconv.ParseFloat(tag, 64)
			}
			fieldValue.SetFloat(r)
		case reflect.String:
			if tag == "" {
				fieldValue.SetString(mock.RandString())
			} else {
				fieldValue.SetString(tag)
			}
		case reflect.Ptr:
			fieldTyp := fieldValue.Type()
			if _, ok := typMap[fieldTyp]; ok {
				return
			}
			if tag == "true" {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldTyp.Elem()))
				}
				autowired(fieldValue, typMap)
			}
		case reflect.Struct:
			if tag == "true" {
				autowired(value.Addr(), typMap)
			}
		}
	}

}
