package initialize

import "io"

type NeedInit interface {
	Init()
}

type Config = NeedInit

type ConfigPlaceholder struct {
}

func (c *ConfigPlaceholder) Init() {
}

type Dao = NeedInit

type DaoPlaceholder struct {
}

func (d *DaoPlaceholder) Init() {
}

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
