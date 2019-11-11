package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

func SupportPrivateFields()  {
	extra.SupportPrivateFields()
}

type StringJson string

func (s *StringJson) MarshalJSON() ([]byte, error) {
	return strings2.ToBytes(string(*s)), nil
}