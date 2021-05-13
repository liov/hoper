package binding

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	stringsi "github.com/liov/hoper/v2/utils/strings"
)

type paramSource gin.Params

var _ setter = paramSource(nil)

func (param paramSource) Peek(key string) ([]string, bool) {
	for i := range param {
		if param[i].Key == key {
			return []string{param[i].Value}, true
		}
	}
	return nil, false
}

// TrySet tries to set a value by request's form source (like map[string][]string)
func (param paramSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByKV(value, field, param, tagValue, opt)
}

type Ctx fiber.Ctx

func (c *Ctx) Peek(key string) ([]string, bool) {
	ctx := (*fiber.Ctx)(c)
	v := stringsi.ToString(ctx.Request().URI().QueryArgs().Peek(key))
	if v != "" {
		return []string{v}, true
	}
	v = ctx.Params(key)
	return []string{v}, v != ""
}
