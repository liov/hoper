package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	dbi "github.com/liov/hoper/go/v2/utils/dao/db"
	"gorm.io/gorm"
)

func (d *contentDao) GetMomentListDB(db *gorm.DB,req *content.MomentListReq) (int64,[]*content.Moment,error) {
	var moments []*content.Moment
	db = db.Table(model.MomentTableName).Where( dbi.PostgreNotDeleted)
	var count int64
	db.Count(&count)
	err:=db.Limit(int(req.PageSize)).Offset(int((req.PageNo - 1) * req.PageSize)).
		Find(&moments).Error
	if err != nil {
		return 0,nil, d.ErrorLog(errorcode.DBError, err,"GetMomentListDB")
	}
	return count,moments,nil
}

func (d *contentDao) GetTopMomentsRedis(conn redis.Cmdable,key string, pageNo int, PageSize int) ([]content.Moment, error) {
	var moments []content.Moment
	exist,err:=conn.Exists(d.Context,key).Result()
	if err != nil {
		return nil, d.ErrorLog(errorcode.RedisErr, err,"GetTopMomentsRedis")
	}
	if exist == 0 {
		return nil, d.ErrorLog(errorcode.DataLoss, err,"GetTopMomentsRedis")
	}
	conn.Pipelined(d.Context, func(pipe redis.Pipeliner) error {
		return nil
	})

	return moments,d.ErrorLog(errorcode.RedisErr, err,"GetTopMomentsRedis")
}
