package initialize

import "io"

type NeedInit interface {
	Init()
}

type Config = NeedInit

type NeedInitPlaceholder struct {
}

type Dao = NeedInit

type DaoField interface {
	Config() Generate
	SetEntity(any)
}

type DaoFieldCloser = io.Closer
type DaoFieldCloser1 interface {
	Close()
}

type Generate interface {
	Generate() any
}
