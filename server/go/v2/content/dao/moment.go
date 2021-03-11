package dao

import (
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	dbi "github.com/liov/hoper/go/v2/utils/dao/db"
	"gorm.io/gorm"
)

func (d *contentDao) GetMomentListDB(db *gorm.DB,req *content.MomentListReq) ([]*content.Moment,error) {
	var moments []*content.Moment
	err := db.Table(model.MomentTableName).Where(`deleted_at = ?`, dbi.PostgreZeroTime).
		Limit(int(req.PageSize)).Offset(int((req.PageNo - 1) * req.PageSize)).
		Find(&moments).Error
	if err != nil {
		return nil, err
	}
	return moments,nil
}

func (d *contentDao) DeleteMomentDB(db *gorm.DB,id,userId uint64) error {
	return db.Table(model.MomentTableName).Where(`id = ? AND user_id = ? `, id,userId).
		UpdateColumns(dbi.DeleteAt(d.TimeString)).Error
}


func (d *contentDao) GetTopMomentsRedis(conn redis.Cmdable,key string, pageNo int, PageSize int) ([]content.Moment, int64, int64) {
	var moments []content.Moment
	exist,err:=conn.Exists(d.Context,key).Result()
	if err != nil {
		d.Error(err.Error())
		return nil, 0, 0
	}
	if exist == 0 {
		return nil, 0, 0
	}
	conn.Pipelined(d.Context, func(pipe redis.Pipeliner) error {

		return nil
	})

	return moments, 0, 0
}
