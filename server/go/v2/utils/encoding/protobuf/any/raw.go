package any

import (
	"github.com/gogo/protobuf/jsonpb"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

type RawJson []byte

func NewAny(v interface{}) RawJson {
	data, err := json.Standard.Marshal(v)
	if err != nil {
		log.Error(err)
	}
	return data
}

func BytesToJsonAny(b []byte) RawJson {
	b = append([]byte{'"'}, b...)
	return append(b, '"')
}

func StringToJsonAny(s string) RawJson {
	return []byte("\"" + s + "\"")
}

func (a *RawJson) MarshalJSON() ([]byte, error) {
	if len(*a) == 0 {
		return []byte("null"), nil
	}
	return *a, nil
}

func (a *RawJson) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return a.MarshalJSON()
}

func (a *RawJson) Size() int {
	return len(*a)
}

func (a *RawJson) MarshalTo(b []byte) (int, error) {
	return copy(b, *a), nil
}

func (a *RawJson) Unmarshal(b []byte) error {
	*a = b
	return nil
}

func (a *RawJson) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func NewPopulatedRawJson(r randyBytesJson, easy bool) *RawJson {
	if !easy && r.Intn(10) != 0 {
	}
	any := RawJson("{}")
	return &any
}
