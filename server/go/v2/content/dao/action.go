package dao

import (
	sqlstd "database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	dbi "github.com/liov/hoper/go/v2/utils/dao/db"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

func (d *contentDao) ActionCountDB(db *gorm.DB, typ content.ContentType, action content.ActionType, refId uint64, changeCount int) error {

	var expr clause.Expr
	var column string
	switch action {
	case content.ActionLike:
		column = "like_count"
	case content.ActionUnlike:
		column = "unlike_count"
	case content.ActionBrowse:
		column = "browse_count"
	case content.ActionComment:
		column = "comment_count"
	case content.ActionCollect:
		column = "collect_count"
	case content.ActionGiveAction:
		column = "give_count"
	case content.ActionReport:
		column = "report_count"
	}
	symbol := "+"
	if changeCount < 0 {
		symbol = "-"
		changeCount = -changeCount
	}
	expr = gorm.Expr(column + symbol + strconv.Itoa(changeCount))

	err := db.Table(model.ContentExtTableName).Where(`type = ? AND ref_id = ?`, typ, refId).
		Update(column, expr).Error
	if err != nil {
		return d.ErrorLog(errorcode.DBError, err, "ActionCountDB")
	}
	return nil
}

func (d *contentDao) HotCountRedis(conn redis.Cmdable, typ content.ContentType, refId uint64, changeCount float64) error {
	key := content.ContentType_name[int32(typ)][7:] + redisi.SortSet
	err := conn.ZIncrBy(d.Context, key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return d.ErrorLog(errorcode.RedisErr, err, "HotCountRedis")
	}
	return nil
}

func (d *contentDao) ActionCountRedis(conn redis.Cmdable, typ content.ContentType, action content.ActionType, refId uint64, changeCount float64) error {
	key := content.ContentType_name[int32(typ)][7:] + content.ActionType_name[int32(action)][6:] + redisi.SortSet
	err := conn.ZIncrBy(d.Context, key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return d.ErrorLog(errorcode.RedisErr, err, "ActionCountRedis")
	}
	return nil
}

// 是否存在全部改成这种形式
func (d *contentDao) CollectIdDB(db *gorm.DB, typ content.ContentType, refId, userId, favId uint64) (uint64, error) {

	// 性能优化之分开写
	sql := `SELECT id FROM "` + model.CollectTableName + `" 
WHERE type = ? AND ref_id = ? AND user_id = ? AND fav_id = ? AND ` + dbi.PostgreNotDeleted + ``
	var id uint64
	err := db.Raw(sql, typ, refId, userId, favId).Row().Scan(&id)
	if err != nil {
		return 0, d.ErrorLog(errorcode.DBError, err, "CollectIdDB")
	}
	return id, nil
}

func (d *contentDao) LikeIdDB(db *gorm.DB, typ content.ContentType, action content.ActionType, refId, userId uint64) (uint64, error) {

	// 性能优化之分开写
	sql := `SELECT id FROM "` + model.ActionTableName(action) + `" 
WHERE type = ? AND ref_id = ? AND action = ? AND user_id = ? AND ` + dbi.PostgreNotDeleted
	var id uint64
	err := db.Raw(sql, typ, refId, action, userId).Row().Scan(&id)
	if err != nil && err != sqlstd.ErrNoRows {
		return 0, d.ErrorLog(errorcode.DBError, err, "LikeIdDB")
	}
	return id, nil
}

func (d *contentDao) DelActionDB(db *gorm.DB, typ content.ContentType, action content.ActionType, refId, userId uint64) error {
	sql := `Update "` + model.ActionTableName(action) + `" SET deleted_at = ?
WHERE type = ? AND ref_id = ? AND action = ? AND user_id = ? AND ` + dbi.PostgreNotDeleted
	err := db.Exec(sql, d.TimeString, typ, refId, action, userId).Error
	if err != nil {
		return d.ErrorLog(errorcode.DBError, err, "DelActionDB")
	}
	return nil
}

func (d *contentDao) DelDB(db *gorm.DB, tableName string, id uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = ?
WHERE id = ? AND ` + dbi.PostgreNotDeleted
	err := db.Exec(sql, d.TimeString, id).Error
	if err != nil {
		return d.ErrorLog(errorcode.DBError, err, "DelDB")
	}
	return nil
}

func (d *contentDao) DelByAuthDB(db *gorm.DB, tableName string, id, userId uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = ?
WHERE id = ?  AND user_id = ? AND ` + dbi.PostgreNotDeleted
	err := db.Exec(sql, d.TimeString, id, userId).Error
	if err != nil {
		return d.ErrorLog(errorcode.DBError, err, "DelByAuthDB")
	}
	return nil
}

func (d *contentDao) ExistsByAuthDB(db *gorm.DB, tableName string, id, userId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE id = ?  AND user_id = ? AND ` + dbi.PostgreNotDeleted + ` LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, userId).Row().Scan(&exists)
	if err != nil {
		return false, d.ErrorLog(errorcode.DBError, err, "ExistsByAuthDB")
	}
	return exists, nil
}

func (d *contentDao) ContainerExistsDB(db *gorm.DB, typ content.ContainerType, id, userId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + model.ContainerTableName + `" 
WHERE id = ?  AND type = ? AND user_id = ? AND ` + dbi.PostgreNotDeleted + ` LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, typ, userId).Row().Scan(&exists)
	if err != nil {
		return false, d.ErrorLog(errorcode.DBError, err, "ContainerExistsDB")
	}
	return exists, nil
}

func (d *contentDao) GetContentActionDB(db *gorm.DB,action content.ActionType,typ content.ContentType,refIds []uint64,userId uint64) ([]model.ContentAction,error) {
	var actions []model.ContentAction
	err := db.Select("ref_id,id AS like_id,action").Table(model.ActionTableName(action)).
		Where("type = ? AND ref_id IN (?) AND user_id = ? AND "+dbi.PostgreNotDeleted,
			typ,refIds,userId).Scan(&actions).Error
	if err != nil {
		return nil, d.ErrorLog(errorcode.DBError,err,"GetContentActionDB")
	}
	return actions,nil
}