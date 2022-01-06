package gormi

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"gorm.io/gorm"
	"strings"
)

func DelDB(db *gorm.DB, tableName string, id uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = now()
WHERE id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `'`
	return db.Exec(sql, id).Error
}

func DelByAuthDB(db *gorm.DB, tableName string, id, userId uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = now()
WHERE id = ?  AND user_id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `'`
	return db.Exec(sql, id, userId).Error
}

func ExistsDB(db *gorm.DB, tableName string, id uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `' LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ExistsByAuthDB(db *gorm.DB, tableName string, id, userId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE id = ?  AND user_id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `' LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, userId).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type Clause []func(db *gorm.DB) *gorm.DB

// db.Scope(ById(1),ByName("a")).First(v)
func (c Clause) ById(id int) Clause {
	return append(c, func(db *gorm.DB) *gorm.DB {
		return db.Where(`id = ?`, id)
	})
}

func (c Clause) ByName(name string) Clause {
	return append(c, func(db *gorm.DB) *gorm.DB {
		return db.Where(`name = ?`, name)
	})
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
