package protobuf

import "github.com/liov/hoper/go/v2/utils/strings2"

type StringJson string

func (s *StringJson) MarshalJSON() ([]byte, error) {
return strings2.ToBytes(string(*s)), nil
}

func (s *StringJson) Size() int {
	return len(*s)
}

func (s *StringJson) MarshalTo(b []byte) (int,error) {
	return copy(b, *s),nil
}

func (s *StringJson) Unmarshal(b []byte) error {
	*s = StringJson(b)
	return nil
}