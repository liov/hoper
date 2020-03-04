package any

import (
	"github.com/gogo/protobuf/jsonpb"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

type BytesJson []byte

func NewAny(v interface{}) BytesJson {
	data, err := json.Json.Marshal(v)
	if err != nil {
		log.Error(err)
	}
	return data
}

func BytesToJsonAny(b []byte) BytesJson {
	b = append([]byte{'"'}, b...)
	return append(b, '"')
}

func StringToJsonAny(s string) BytesJson {
	return []byte("\"" + s + "\"")
}

func (a *BytesJson) MarshalJSON() ([]byte, error) {
	if len(*a) == 0 {
		return []byte("null"), nil
	}
	return *a, nil
}

func (a *BytesJson) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return a.MarshalJSON()
}

func (a *BytesJson) Size() int {
	return len(*a)
}

func (a *BytesJson) MarshalTo(b []byte) (int, error) {
	return copy(b, *a), nil
}

func (a *BytesJson) Unmarshal(b []byte) error {
	*a = b
	return nil
}

func (a *BytesJson) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i -= len(*a)
	copy(dAtA[i:], *a)
	return len(*a), nil
}

type randyBytesJson interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func NewPopulatedBytesJson(r randyBytesJson, easy bool) *BytesJson {
	if !easy && r.Intn(10) != 0 {
	}
	any := BytesJson("{}")
	return &any
}
