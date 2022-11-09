package errorsi

type Unwrap interface {
	Unwrap(err error) error
}

type Is interface {
	Is(err error) bool
}

type WarpError struct {
	ErrRep
	err error
}

func (x *WarpError) Error() string {
	return x.Message
}

func (x *WarpError) Unwrap() error {
	return x.err
}
