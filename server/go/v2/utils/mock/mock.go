package mock

import (
	"math/rand"
	"reflect"
	"time"
)

//
func init() {
	rand.Seed(time.Now().UnixNano())
}

func Mock(v interface{}) {
	value := reflect.ValueOf(v)
	typMap := make(map[reflect.Type]int)
	mock(value, typMap)
}

//数组长度
const length = 1

//一个类型最大重复次数
const times = 3

func mock(value reflect.Value, typMap map[reflect.Type]int) {
	typ := value.Type()
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := uint64(rand.Int63n(10000))
		value.SetUint(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := rand.Int63n(10000)
		value.SetInt(r)
	case reflect.Float32, reflect.Float64:
		r := rand.ExpFloat64()
		value.SetFloat(r)
	case reflect.String:
		value.SetString(RandString())
	case reflect.Ptr:
		if count := typMap[typ]; count == times {
			return
		}
		typMap[typ] = typMap[typ] + 1
		if value.IsNil() && value.CanSet() {
			value.Set(reflect.New(typ.Elem()))
		}else {
			return
		}
		mock(value.Elem(), typMap)
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			mock(field, typMap)
		}
	case reflect.Array:
		for i := 0; i < value.Len(); i++ {
			mock(value.Index(i), typMap)
		}
	case reflect.Slice:
		value.Set(reflect.MakeSlice(typ, length, length))
		for i := 0; i < length; i++ {
			mock(value.Index(i), typMap)
		}
	case reflect.Map:
		value.Set(reflect.MakeMapWithSize(typ, length))
		for i := 0; i < length; i++ {
			mk := reflect.New(typ.Key()).Elem()
			mock(mk, typMap)
			mv := reflect.New(typ.Elem()).Elem()
			mock(mv, typMap)
			value.SetMapIndex(mk, mv)
		}
	}
}

func RandString() string {
	runes := make([]rune, 5)
	for i := range runes {
		runes[i] = randRune()
	}
	return string(runes)
}

func randRune() rune {
	r := rand.Intn(20901)
	return rune(r + 19968)
}
