package _interface

type Get interface {
	Get(key string) string
}

type Init interface {
	Init()
}
