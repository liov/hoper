package conf_center

type ConfigCenter interface {
	SetConfig(func([]byte)) error
}
