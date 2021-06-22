package conf_center

import "github.com/liov/hoper/v2/utils/configor/apollo"

type Apollo struct {
	apollo.Config
}

func (e *Apollo) HandleConfig(handle func([]byte)) error {
	return nil
}
