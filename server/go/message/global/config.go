package global

import (
	"github.com/hopeio/gox/idgen"
	"github.com/hopeio/initialize"
)

type config struct {
	initialize.EmbeddedPresets
	ServerId uint16
}

func (c *config) AfterInject() {
	SF = idgen.NewSnowflake(c.ServerId, 10)
}
