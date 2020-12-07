package configor

type ConfigCenter interface {
	NewClient() Client
	GetConfigCenterService(string) (Service,error)
}

type Client interface {
	Listener(func([]byte))
	GetConfigAllInfoHandle(func([]byte)) error
}

type Service interface {

}
