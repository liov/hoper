package any

import (
	"github.com/liov/hoper/server/go/lib/utils/encoding/json/iterator"
)

func NewAny(v interface{}) (*RawJson, error) {
	data, err := iterator.Standard.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &RawJson{B: data}, nil
}

func BytesToJsonAny(b []byte) *RawJson {
	b = append([]byte{'"'}, b...)
	return &RawJson{B: append(b, '"')}
}

func StringToJsonAny(s string) *RawJson {
	return &RawJson{B: []byte("\"" + s + "\"")}
}

func (a *RawJson) MarshalJSON() ([]byte, error) {
	if len(a.B) == 0 {
		return []byte("null"), nil
	}
	return a.B, nil
}

func (a *RawJson) Size() int {
	return len(a.B)
}

func (a *RawJson) MarshalTo(b []byte) (int, error) {
	return copy(b, a.B), nil
}

func (a *RawJson) Unmarshal(b []byte) error {
	a.B = b
	return nil
}

func (a *RawJson) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i -= len(a.B)
	copy(dAtA[i:], a.B)
	return len(a.B), nil
}
