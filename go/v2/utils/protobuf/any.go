package protobuf

import (
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

type Any []byte

func NewAny(v interface{}) Any {
	data, err := json.Json.Marshal(v)
	if err != nil {
		log.Error(err)
	}
	return data
}

func BytesToJsonAny(b []byte) Any {
	b = append([]byte{'"'}, b...)
	return append(b, '"')
}

func StringToJsonAny(s string) Any {
	return []byte("\"" + s + "\"")
}

func (a *Any) MarshalJSON() ([]byte, error) {
	if len(*a) == 0 {
		return []byte("null"), nil
	}
	return *a, nil
}

func (a *Any) Size() int {
	return len(*a)
}

func (a *Any) MarshalTo(b []byte) (int, error) {
	return copy(b, *a), nil
}

func (a *Any) Unmarshal(b []byte) error {
	*a = b
	return nil
}
