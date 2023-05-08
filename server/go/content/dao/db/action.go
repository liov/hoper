package db

import (
	sqlstd "database/sql"
	"github.com/hopeio/pandora/protobuf/errorcode"
	dbi "github.com/hopeio/pandora/utils/dao/db/const"
	clausei "github.com/hopeio/pandora/utils/dao/db/gorm/clause"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

func (d *ContentDBDao) ActionCount(typ content.ContentType, action content.ActionType, refId uint64, changeCount int) error {
	ctxi := d.Context
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
	case content.ActionGive:
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

	err := d.db.Table(model.ContentExtTableName).Where(`type = ? AND ref_id = ?`, typ, refId).
		Update(column, expr).Error
	if err != nil {
		return ctxi.ErrorLog(errorcode.DBError, err, "ActionCount")
	}
	return nil
}

func (d *ContentDBDao) LikeId(typ content.ContentType, action content.ActionType, refId, userId uint64) (uint64, error) {
	ctxi := d.Context
	// 性能优化之分开写
	sql := `SELECT id FROM "` + model.ActionTableName(action) + `" 
WHERE type = ? AND ref_id = ? AND action = ? AND user_id = ?` + dbi.WithNotDeleted
	var id uint64
	err := d.db.Raw(sql, typ, refId, action, userId).Row().Scan(&id)
	if err != nil && err != sqlstd.ErrNoRows {
		return 0, ctxi.ErrorLog(errorcode.DBError, err, "LikeId")
	}
	return id, nil
}

func (d *ContentDBDao) DelAction(typ content.ContentType, action content.ActionType, refId, userId uint64) error {
	ctxi := d.Context
	sql := `Update "` + model.ActionTableName(action) + `" SET deleted_at = ?
WHERE type = ? AND ref_id = ? AND action = ? AND user_id = ?` + dbi.WithNotDeleted
	err := d.db.Exec(sql, ctxi.TimeString, typ, refId, action, userId).Error
	if err != nil {
		return ctxi.ErrorLog(errorcode.DBError, err, "DelAction")
	}
	return nil
}

func (d *ContentDBDao) Del(tableName string, id uint64) error {
	ctxi := d.Context
	sql := `Update "` + tableName + `" SET deleted_at = ?
WHERE id = ?` + dbi.WithNotDeleted
	err := d.db.Exec(sql, ctxi.TimeString, id).Error
	if err != nil {
		return ctxi.ErrorLog(errorcode.DBError, err, "Del")
	}
	return nil
}

func (d *ContentDBDao) DelByAuth(tableName string, id, userId uint64) error {
	ctxi := d.Context
	sql := `Update "` + tableName + `" SET deleted_at = ?
WHERE id = ?  AND user_id = ?` + dbi.WithNotDeleted
	err := d.db.Exec(sql, ctxi.TimeString, id, userId).Error
	if err != nil {
		return ctxi.ErrorLog(errorcode.DBError, err, "DelByAuth")
	}
	return nil
}

func (d *ContentDBDao) ExistsByAuth(tableName string, id, userId uint64) (bool, error) {
	ctxi := d.Context
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE id = ?  AND user_id = ?` + dbi.WithNotDeleted + ` LIMIT 1)`
	var exists bool
	err := d.db.Raw(sql, id, userId).Scan(&exists).Error
	if err != nil {
		return false, ctxi.ErrorLog(errorcode.DBError, err, "ExistsByAuth")
	}
	return exists, nil
}

func (d *ContentDBDao) ContainerExists(typ content.ContainerType, id, userId uint64) (bool, error) {
	ctxi := d.Context
	sql := `SELECT EXISTS(SELECT * FROM "` + model.ContainerTableName + `" 
WHERE id = ?  AND type = ? AND user_id = ?` + dbi.WithNotDeleted + ` LIMIT 1)`
	var exists bool
	err := d.db.Raw(sql, id, typ, userId).Scan(&exists).Error
	if err != nil {
		return false, ctxi.ErrorLog(errorcode.DBError, err, "ContainerExists")
	}
	return exists, nil
}

func (d *ContentDBDao) GetContentActions(action content.ActionType, typ content.ContentType, refIds []uint64, userId uint64) ([]model.ContentAction, error) {
	ctxi := d.Context
	var actions []model.ContentAction
	err := d.db.Select("id,ref_id,action").Table(model.ActionTableName(action)).
		Where("type = ? AND ref_id IN (?) AND user_id = ?"+dbi.WithNotDeleted,
			typ, refIds, userId).Scan(&actions).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetContentActions")
	}
	return actions, nil
}

func (d *ContentDBDao) GetLike(likeId, userId uint64) (*model.ContentAction, error) {
	ctxi := d.Context
	var action model.ContentAction
	err := d.db.Select("id,ref_id,action,type").Table(model.LikeTableName).
		Where("id = ? AND user_id = ?"+dbi.WithNotDeleted,
			likeId, userId).Scan(&action).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetContentActions")
	}
	return &action, nil
}

func (d *ContentDBDao) GetCollects(typ content.ContentType, refIds []uint64, userId uint64) ([]model.ContentCollect, error) {
	ctxi := d.Context
	var collects []model.ContentCollect
	err := d.db.Select("id,ref_id,fav_id").Table(model.CollectTableName).
		Where("type = ? AND ref_id IN (?) AND user_id = ?"+dbi.WithNotDeleted,
			typ, refIds, userId).Scan(&collects).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetContentActions")
	}
	return collects, nil
}

func (d *ContentDBDao) GetComments(typ content.ContentType, refId, rootId uint64, pageNo, pageSize int) (int64, []*content.Comment, error) {
	ctxi := d.Context
	db := d.db.Table(model.CommentTableName).Where(`type = ? AND ref_id = ? AND root_id = ?`+dbi.WithNotDeleted, typ, refId, rootId)
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "Find")
	}
	var clauses []clause.Expression
	clauses = append(clauses, clausei.Page(pageNo, pageSize))
	var comments []*content.Comment
	err = db.Clauses(clauses...).Find(&comments).Error
	if err != nil {
		return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "Find")
	}
	return count, comments, nil
}
