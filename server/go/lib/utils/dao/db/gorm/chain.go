package gormi

import (
	"context"
	"gorm.io/gorm"
)

type ChainDao struct {
	Ctx          context.Context
	DB, OriginDB *gorm.DB
}

func NewChainDao(ctx context.Context, db *gorm.DB) *ChainDao {
	return &ChainDao{
		ctx, db, db,
	}
}

func (c *ChainDao) ResetDB() {
	c.DB = c.OriginDB
}

func (c *ChainDao) ById(id int) *ChainDao {
	c.DB.Where(`id = ?`, id)
	return c
}

func (c *ChainDao) ByName(name string) *ChainDao {
	c.DB.Where(`name = ?`, name)
	return c
}
