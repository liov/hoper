package clause

import (
	"context"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	"gorm.io/gorm"
	"strings"
)

type ChainDB func(db *gorm.DB) *gorm.DB

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

func NewChainDB(ctx context.Context) ChainDB {
	return func(db *gorm.DB) *gorm.DB {
		return db.WithContext(ctx)
	}
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
