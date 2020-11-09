package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

func SupportPrivateFields() {
	extra.SupportPrivateFields()
}

func Marshal(v interface{}) ([]byte, error) {
	return Json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return Json.Unmarshal(data, v)
}
