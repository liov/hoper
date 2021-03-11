package dao

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *contentDao) ActionCountDB(db *gorm.DB,typ content.ContentType, action content.ActionType, refId uint64,changeCount int) error {

	var expr clause.Expr
	var column string
	switch action {
	case content.ActionLike, content.ActionUnlike:
		column = "like_count"
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
	if action == content.ActionUnlike {
		symbol = "-"
	}
	expr = gorm.Expr(column + symbol + strconv.Itoa(changeCount))

	return db.Table(model.ContentExtTableName).Where(`type = ? AND ref_id = ?`,typ,refId).
		Update(column, expr).Error
}

func (d *contentDao) LikeCountRedis(conn redis.Cmdable, typ content.ContentType,refId uint64, changeCount float64) error {
	key:=content.ContentType_name[int32(typ)][7:] + "_Sorted_Set"
	return conn.ZIncrBy(d.Context,key,changeCount, strconv.FormatUint(refId,10)).Err()
}