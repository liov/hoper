package valid

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/liov/hoper/go/v2/utils/log"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var (
	Validate *validator.Validate
	trans    ut.Translator
)

func init() {
	zh := zh.New()
	uni := ut.New(zh, zh)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("zh")

	Validate = validator.New()
	zh_translations.RegisterDefaultTranslations(Validate, trans)
	Validate.RegisterTagNameFunc(func(sf reflect.StructField) string {
		return sf.Tag.Get("comment")
	})
	Validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString(`^1[0-9]{10}$`, fl.Field().String())
		return match
	})
	Validate.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0}必须是一个有效的手机号!", true)
	}, translateFunc)
}

func Trans(err error) string {
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
