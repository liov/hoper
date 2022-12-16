package initialize

type Config[D any] interface {
	Build() (*D, func() error)
}

type Dao[C Config[D], D any] struct {
	Conf   C
	Client *D
	Close  func() error
}

func (d *Dao[C, D]) SetClient() {
	d.Client, d.Close = d.Conf.Build()
}
