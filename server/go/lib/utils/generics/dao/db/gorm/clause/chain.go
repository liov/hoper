package clausei

import (
	"context"
	"gorm.io/gorm"
)

type ChainDao[C context.Context] struct {
	Ctx          C
	DB, OriginDB *gorm.DB
}

func (c *ChainDao[C]) ResetDB() {
	c.DB = c.OriginDB
}

func (c *ChainDao[C]) ById(id int) *ChainDao[C] {
	c.DB.Where(`id = ?`, id)
	return c
}

func (c *ChainDao[C]) ByName(name string) *ChainDao[C] {
	c.DB.Where(`name = ?`, name)
	return c
}
