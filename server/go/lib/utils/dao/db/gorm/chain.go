package gormi

import (
	"context"
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"gorm.io/gorm"
	"strings"
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

type Clause []func(db *gorm.DB) *gorm.DB

type Expression dbi.Expression

func (e *Expression) Clause() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(e.Field+(*dbi.Expression)(e).Operation.SQL(), e.Value...)
	}
}

func NewScope(field string, op dbi.Operation, args ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+op.SQL(), args...)
	}
}

// db.Scope(ById(1),ByName("a")).First(v)
func (c Clause) ById(id int) Clause {
	return append(c, NewScope("id", dbi.Equal, id))
}

func (c Clause) ByName(name string) Clause {
	return append(c, func(db *gorm.DB) *gorm.DB {
		return db.Where(`name = ?`, name)
	})
}

func (c Clause) Exec(db *gorm.DB) *gorm.DB {
	db = db.Scopes(c...)
	return db
}

type Clause2 struct {
	Expr []string
	Var  []interface{}
}

// db.Scope(ById(1),ByName("a").Build()).First(v)
func (c *Clause2) ById(id int) *Clause2 {
	c.Expr = append(c.Expr, `id = ?`)
	c.Var = append(c.Var, id)
	return c
}

func (c *Clause2) ByName(name string) *Clause2 {
	c.Expr = append(c.Expr, `name = ?`)
	c.Var = append(c.Var, name)
	return c
}

func (c *Clause2) Build(db *gorm.DB) *gorm.DB {
	db = db.Where(strings.Join(c.Expr, " AND "), c.Var...)
	return db
}

func (c *Clause2) Exec(db *gorm.DB) *gorm.DB {
	db = db.Scopes(c.Build)
	return db
}
