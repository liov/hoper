package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/liov/hoper/server/go/lib/utils/log"
)

var (
	Validator *validator.Validate
	trans     ut.Translator
)

func init() {
	zhcn := zh.New()
	uni := ut.New(zhcn, zhcn)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("zh")

	Validator = validator.New()
	zh_translations.RegisterDefaultTranslations(Validator, trans)
	Validator.RegisterTagNameFunc(func(sf reflect.StructField) string {
		if annotation := sf.Tag.Get("annotation"); annotation != "" {
			return annotation
		}
		if json := sf.Tag.Get("json"); json != "" {
			return json
		}
		return sf.Name
	})
	Validator.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^1[0-9]{10}$`, fl.Field().String())
		return match
	})
	Validator.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0}必须是一个有效的手机号!", true)
	}, translateFunc)
}

func Trans(err error) string {
	if err == nil {
		return ""
	}
	var msg []string
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		log.Error("无效的参数")
		return ""
	}
	for _, v := range ve.Translate(trans) {
		msg = append(msg, v)
	}
	return strings.Join(msg, ",")
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Errorf("警告: 翻译字段错误: %#v", fe)
		return fe.(error).Error()
	}

	return t
}
