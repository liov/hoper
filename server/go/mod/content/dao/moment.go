package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	gormi "github.com/liov/hoper/server/go/lib/utils/dao/db/gorm"
	"github.com/liov/hoper/server/go/mod/content/model"
	"github.com/liov/hoper/server/go/mod/protobuf/content"
	"gorm.io/gorm/clause"
)

func (d *contentDao) GetMomentListDB(req *content.MomentListReq) (int64, []*content.Moment, error) {
	ctxi := d.Ctx
	var moments []*content.Moment
	db := d.Table(model.MomentTableName).Where(dbi.PostgreNotDeleted)
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "Count")
	}
	var clauses []clause.Expression
	clauses = append(clauses, gormi.Page(int(req.PageNo), int(req.PageSize)))
	err = db.Clauses(clauses...).Order("created_at desc").Find(&moments).Error
	if err != nil {
		return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "GetMomentListDB")
	}
	return count, moments, nil
}

func (d *contentDao) GetTopMomentsRedis(conn redis.Cmdable, key string, pageNo int, PageSize int) ([]content.Moment, error) {
	ctxi := d.Ctx
	var moments []content.Moment
	exist, err := conn.Exists(ctxi.Context, key).Result()
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr, err, "GetTopMomentsRedis")
	}
	if exist == 0 {
		return nil, ctxi.ErrorLog(errorcode.DataLoss, err, "GetTopMomentsRedis")
	}
	conn.Pipelined(ctxi.Context, func(pipe redis.Pipeliner) error {
		return nil
	})

	return moments, ctxi.ErrorLog(errorcode.RedisErr, err, "GetTopMomentsRedis")
}
