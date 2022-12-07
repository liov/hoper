package iterator

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var Standard = jsoniter.ConfigCompatibleWithStandardLibrary

func SupportPrivateFields() {
	extra.SupportPrivateFields()
}

func Marshal(v interface{}) ([]byte, error) {
	return Standard.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return Standard.Unmarshal(data, v)
}

var WithPrivateField = jsoniter.Config{
	IndentionStep:                 4,
	MarshalFloatWith6Digits:       true,
	EscapeHTML:                    true,
	SortMapKeys:                   true,
	UseNumber:                     true,
	ObjectFieldMustBeSimpleString: true,
}.Froze()
