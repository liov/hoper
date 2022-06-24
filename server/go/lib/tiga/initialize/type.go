package initialize

type NeedInit interface {
	Init()
}

type Config interface {
	NeedInit
}

type Dao interface {
	Close()
	NeedInit
}

type DaoField interface {
	Config() Generate
	SetEntity(any)
}

type Generate interface {
	Generate() any
}
