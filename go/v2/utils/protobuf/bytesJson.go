package protobuf

import (
	"github.com/gogo/protobuf/jsonpb"
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

func (a *Any) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error)  {
	return a.MarshalJSON()
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

func (a *Any) MarshalToSizedBuffer(dAtA []byte) (int,error) {
	i := len(dAtA)
	i -= len(*a)
	copy(dAtA[i:], *a)
	return len(*a),nil
}

type randyBytesJson interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func NewPopulatedAny(r randyBytesJson,easy bool) *Any {
	if !easy && r.Intn(10) != 0 {
	}
	any := Any("{}")
	return &any
}