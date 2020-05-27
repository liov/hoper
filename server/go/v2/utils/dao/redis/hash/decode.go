package hash

import (
	"reflect"
	"strconv"
)

type decodeState struct {
	strings []string
}

/*func (d *decodeState) decode(index int,v reflect.Value)  {
	v := reflect.ValueOf(v)
	t := v.Type()
	switch t.Kind() {
	case reflect.Bool:
		e.strings = append(e.strings, t.Name(), v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Ptr:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}
*/

func UnMarshal(v interface{}, redisArgs []string) {
	uValue := reflect.ValueOf(v).Elem()
	for i := 0; i < len(redisArgs); i += 2 {
		fieldValue := uValue.FieldByName(redisArgs[i])
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, _ := strconv.ParseInt(redisArgs[i+1], 10, 64)
			fieldValue.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, _ := strconv.ParseUint(redisArgs[i+1], 10, 64)
			fieldValue.SetUint(v)
		case reflect.String:
			fieldValue.SetString(redisArgs[i+1])
		case reflect.Float32, reflect.Float64:
			v, _ := strconv.ParseFloat(redisArgs[i+1], 64)
			fieldValue.SetFloat(v)
		case reflect.Bool:
			v, _ := strconv.ParseBool(redisArgs[i+1])
			fieldValue.SetBool(v)
		}
	}
}
