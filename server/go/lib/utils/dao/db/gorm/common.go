package gormi

import (
	contexti "github.com/liov/hoper/server/go/lib/utils/context"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	"gorm.io/gorm"
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

type CommonChainDao struct {
	Ctx *contexti.RequestContext
	*gorm.DB
	originDB *gorm.DB
}

func (c *CommonChainDao) ResetDB() {
	c.DB = c.originDB
}

func New(ctx *contexti.RequestContext, db *gorm.DB) *CommonChainDao {
	db = db.Session(&gorm.Session{Context: contexti.SetTranceId(ctx.TraceID), NewDB: true})
	return &CommonChainDao{Ctx: ctx, DB: db, originDB: db}
}

func (c *CommonChainDao) NewDB(db *gorm.DB) *CommonChainDao {
	return &CommonChainDao{c.Ctx, db, db}
}
