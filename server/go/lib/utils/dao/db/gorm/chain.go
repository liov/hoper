package gormi

import (
	"context"
	"gorm.io/gorm"
)

type ChainDao struct {
	DB, OriginDB *gorm.DB
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

type ChainDB func(db *gorm.DB) *gorm.DB

func ChainDBHelper() {

}

func (c ChainDB) ById(id int) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return c(db).Where("id = ?", id)
	}
}

func (c ChainDB) ByName(name string) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return c(db).Where("name = ?", name)
	}
}

func (c ChainDB) List(db *gorm.DB) {
	c(db).Find(nil)
}

type TestChainDBDao struct {
	ChainDB
}

func NewChainDB(ctx context.Context) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return db.WithContext(ctx)
	}
}

func NewTestChainDBDao(ctx context.Context) *TestChainDBDao {
	return &TestChainDBDao{NewChainDB(ctx)}
}

func testChainDBDao() {
	dao := NewTestChainDBDao(context.Background())
	db := new(gorm.DB)
	dao.ById(1).ByName("a").List(db)
}
