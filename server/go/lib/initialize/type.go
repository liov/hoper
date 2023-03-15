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
	Config() any
	SetEntity()
	io.Closer
}

type DaoFieldCloser = io.Closer
type DaoFieldCloser1 interface {
	Close()
}

type CloseFunc func() error

// Deprecated
type Generate interface {
	Generate() any
}

// TODO:泛型统一的Dao结构
type DaoConfig[D any] interface {
	Build() (*D, CloseFunc)
}

type DaoEntity[C DaoConfig[D], D any] struct {
	Conf   C
	Client *D
	close  CloseFunc
}

func (d *DaoEntity[C, D]) Config() any {
	return d.Conf
}

func (d *DaoEntity[C, D]) SetEntity() {
	d.Client, d.close = d.Conf.Build()
}

func (d *DaoEntity[C, D]) Close() error {
	if d.Close != nil {
		return d.close()
	}
	return nil
}
