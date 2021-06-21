package conf_center

type Apollo struct {
}

func (e *Apollo) SetConfig(handle func([]byte)) error {
	return nil
}
