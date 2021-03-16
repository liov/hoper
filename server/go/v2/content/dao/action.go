package dao

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	dbi "github.com/liov/hoper/go/v2/utils/dao/db"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	return db.Table(model.ContentExtTableName).Where(`type = ? AND ref_id = ?`, typ, refId).
		Update(column, expr).Error
}

func (d *contentDao) LikeCountRedis(conn redis.Cmdable, typ content.ContentType, refId uint64, changeCount float64) error {
	key := content.ContentType_name[int32(typ)][7:] + redisi.SortSet
	return conn.ZIncrBy(d.Context, key, changeCount, strconv.FormatUint(refId, 10)).Err()
}

func (d *contentDao) ActionExists(db *gorm.DB, typ content.ContentType, action content.ActionType, refId, userId uint64) (bool, error) {
	var tableName string
	switch action {
	case content.ActionBrowse:
		tableName = model.BrowserTableName
	case content.ActionLike, content.ActionUnlike:
		tableName = model.LikeTableName
	case content.ActionComment:
		tableName = model.CommentTableName
	case content.ActionCollect:
		tableName = model.CollectTableName
	case content.ActionReport:
		tableName = model.ReportTableName
	case content.ActionGiveAction:
		tableName = model.GiveTableName
	case content.ActionApprove:
		tableName = model.ApproveTableName
	}
	// 性能优化之分开写
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE type = ? AND ref_id = ? AND action = ? AND user_id = ? AND deleted_at = '` + dbi.PostgreZeroTime2 + `' LIMIT 1)`
	var exists bool
	err := db.Raw(sql, typ, refId, action,userId).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return false, nil
}
