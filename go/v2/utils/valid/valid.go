package valid

import (
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
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
}

func Trans(err error) string {
	var msg []string
	for _,v := range err.(validator.ValidationErrors).Translate(trans) {
		msg = append(msg, v)
	}
	return strings.Join(msg, ",")
}
