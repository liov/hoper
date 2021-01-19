package def

type Less func(a, b interface{}) bool
type Equal func(a, b interface{}) bool

type CmpKey interface {
	CmpKey() uint64
}