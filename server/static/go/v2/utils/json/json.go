package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

func SupportPrivateFields() {
	extra.SupportPrivateFields()
}
