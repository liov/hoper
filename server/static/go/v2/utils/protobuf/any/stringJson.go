package any

import (
	"github.com/gogo/protobuf/jsonpb"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

type StringJson string

func (s *StringJson) MarshalJSON() ([]byte, error) {
	if *s == "" {
		return []byte("null"), nil
	}
	return strings2.ToBytes(string(*s)), nil
}

func (s *StringJson) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return s.MarshalJSON()
}

func (s *StringJson) Size() int {
	return len(*s)
}

func (s *StringJson) MarshalTo(b []byte) (int, error) {
	return copy(b, *s), nil
}

func (s *StringJson) Unmarshal(b []byte) error {
	*s = StringJson(b)
	return nil
}

func (s *StringJson) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i -= len(*s)
	copy(dAtA[i:], *s)
	return len(*s), nil
}

type randyStringJson interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func NewPopulatedStringJson(r randyStringJson, easy bool) *StringJson {
	if !easy && r.Intn(10) != 0 {
	}
	any := StringJson("{}")
	return &any
}

/*
如果用 protoc 生成的文件，手动实现MarshalJSON
type StringJson struct {
	S string
}

func (p *StringJson) MarshalJSON() ([]byte, error) {
	return strings2.ToBytes(p.S), nil
}
*/
