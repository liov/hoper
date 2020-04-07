package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
)

var handleMap = make(map[string]reflect.Value)

func commonHandler(validateSession func(string, interface{}) error) iris.Handler {
	return func(ctx iris.Context) {
		handle, ok := handleMap[ctx.Path()]
		if !ok {
			ctx.NotFound()
			return
		}
		handleTyp := handle.Type()
		handleNumIn := handleTyp.NumIn()
		if handleNumIn != 0 {
			params := make([]reflect.Value, handleNumIn)
			for i := 0; i < handleNumIn; i++ {
				param := reflect.New(handleTyp.In(i))
				params = append(params, param)
				if handleTyp.In(i).Implements(contextType) {
					const authErr = `{"msg":"未登录"}`
					tokenString := ctx.GetCookie("token")
					if len(tokenString) == 0 {
						tokenString = ctx.GetHeader("Authorization")
					}
					if len(tokenString) == 0 {
						ctx.WriteString(authErr)
						return
					}
					if err := validateSession(tokenString, param); err != nil {
						ctx.WriteString(authErr)
						return
					}
				} else {
					//这块借鉴iris schema有待优化
					if ctx.Method() == http.MethodGet {
						unmarshalParams(ctx.Request().URL.Query(), param)
					} else {

					}
				}
			}
			handle.Call(params)
		}
	}
}

func unmarshalParams(query url.Values, param reflect.Value) error {
	param = param.Elem()
	for i := 0; i < param.Type().NumField(); i++ {
		name := lowerFirst(param.Type().Field(i).Name)
		fieldName := strings.TrimSpace(param.Type().Field(i).Tag.Get("json"))
		defaultValue := strings.TrimSpace(param.Type().Field(i).Tag.Get("default"))
		if fieldName != "" {
			name = fieldName
		}
		var strValue string

		if arr := query[name]; len(arr) == 0 {
			strValue = defaultValue
		} else {
			strValue = arr[0]
		}

		fieldKind := param.Type().Field(i).Type.Kind().String()

		// 对不同类型的数据进行校验
		if strings.Contains(fieldKind, "uint") {
			if strValue == "" {
				param.Field(i).SetUint(0)
				continue
			}
			value, err := strconv.ParseUint(strValue, 10, 64)
			if err != nil {
				return errors.New(fmt.Sprintf("%s应为uint", name))
			}
			param.Field(i).SetUint(value)
		} else if strings.Contains(fieldKind, "int") {
			if strValue == "" {
				param.Field(i).SetInt(0)
				continue
			}
			value, err := strconv.ParseInt(strValue, 10, 64)
			if err != nil {
				return errors.New(fmt.Sprintf("%s应为int", name))
			}
			param.Field(i).SetInt(value)
		} else if strings.Contains(fieldKind, "string") {
			param.Field(i).SetString(strValue)
		} else if strings.Contains(fieldKind, "float") {
			if strValue == "" {
				param.Field(i).SetFloat(0.0)
				continue
			}
			value, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return errors.New(fmt.Sprintf("%s应为float", name))
			}
			param.Field(i).SetFloat(value)
		} else if strings.Contains(fieldKind, "bool") {
			if strValue == "" {
				param.Field(i).SetBool(false)
				continue
			}
			value, err := strconv.ParseBool(strValue)
			if err != nil {
				return errors.New(fmt.Sprintf("%s应为bool", name))
			}
			param.Field(i).SetBool(value)
		}
	}
	return nil
}
