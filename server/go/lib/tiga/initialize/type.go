package initialize

type NeedInit interface {
	Init()
}

type Config interface {
	NeedInit
}

type ConfigPlaceholder struct {
}

func (c *ConfigPlaceholder) Init() {
}

type Dao interface {
	Close()
	NeedInit
}

type DaoPlaceholder struct {
}

func (d *DaoPlaceholder) Init() {
}
func (d *DaoPlaceholder) Close() {
}

type DaoField interface {
	Config() Generate
	SetEntity(any)
}

type Generate interface {
	Generate() any
}
