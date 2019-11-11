package protobuf

type Any []byte

func (a *Any) MarshalJSON() ([]byte, error) {
	return *a, nil
}

func (a *Any) Size() int {
	return len(*a)
}

func (a *Any) MarshalTo(b []byte) (int,error) {
	return copy(b, *a),nil
}

func (a *Any) Unmarshal(b []byte) error {
	*a = b
	return nil
}