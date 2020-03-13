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
	value := reflect.ValueOf(v).Elem()
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r := rand.Uint64()
		value.SetUint(r)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r := rand.Int63()
		value.SetInt(r)
	case reflect.Float32, reflect.Float64:
		r := rand.Float64()
		value.SetFloat(r)
	case reflect.String:
		value.SetString(RandString())
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		Mock(value.Interface())
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			Mock(field.Addr().Interface())
		}
	case reflect.Array:

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
