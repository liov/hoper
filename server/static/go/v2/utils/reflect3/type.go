package reflect3

import (
	"reflect"
	"strconv"
)

func GetDereferenceType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Ptr:
		return GetDereferenceType(typ.Elem())
	case reflect.Slice:
		return GetDereferenceType(typ.Elem())
	case reflect.Array:
		return GetDereferenceType(typ.Elem())
	case reflect.Chan:
		return GetDereferenceType(typ.Elem())
	case reflect.Map:
		return GetDereferenceType(typ.Elem())
	}
	return typ
}

type StructTag string

func GetCustomizeTag(customize, key string) string {
	v, _ := StructTag(customize).Lookup(key)
	return v
}

func (tag StructTag) Lookup(key string) (value string, ok bool) {

	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '\'' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '\'' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '\'' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
