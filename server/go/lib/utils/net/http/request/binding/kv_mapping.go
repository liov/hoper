package binding

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/valyala/fasthttp"
)

type Arg interface {
	Peek(key string) ([]string, bool)
}

type Args []Arg

func (args Args) Peek(key string) (v []string, ok bool) {
	for i := range args {
		if v, ok = args[i].Peek(key); ok {
			return
		}
	}
	return
}

func (args Args) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSetted bool, err error) {
	return setByKV(value, field, args, key, opt)
}

type argsSource fasthttp.Args

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form *argsSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByKV(value, field, form, tagValue, opt)
}

func (form *argsSource) Peek(key string) ([]string, bool) {
	v := stringsi.ToString((*fasthttp.Args)(form).Peek(key))
	return []string{v}, v != ""
}

func setByKV(value reflect.Value, field reflect.StructField, kv Arg, tagValue string, opt setOptions) (isSetted bool, err error) {
	vs, ok := kv.Peek(tagValue)
	if !ok && !opt.isDefaultExists {
		return false, nil
	}

	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		return true, setSlice(vs, value, field)
	case reflect.Array:
		if !ok {
			vs = []string{opt.defaultValue}
		}
		if len(vs) != value.Len() {
			return false, fmt.Errorf("%q is not valid value for %s", vs, value.Type().String())
		}
		return true, setArray(vs, value, field)
	default:
		var val string
		if !ok {
			val = opt.defaultValue
		}

		if len(vs) > 0 {
			val = vs[0]
		}
		return true, setWithProperType(val, value, field)
	}
}

type ctxSource fiber.Ctx

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form *ctxSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByKV(value, field, form, tagValue, opt)
}

func (form *ctxSource) Peek(key string) ([]string, bool) {
	v := (*fiber.Ctx)(form).Params(key)
	return []string{v}, v != ""
}

type reqSource fasthttp.RequestHeader

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form *reqSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByKV(value, field, form, tagValue, opt)
}

func (form *reqSource) Peek(key string) ([]string, bool) {
	v := stringsi.ToString((*fasthttp.RequestHeader)(form).Peek(key))
	return []string{v}, v != ""
}
