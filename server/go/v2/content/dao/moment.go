package dao

import (
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	dbi "github.com/liov/hoper/go/v2/utils/dao/db"
	"github.com/liov/hoper/go/v2/utils/log"
	"gorm.io/gorm"
)

type contentDao struct {
	*user.Ctx
}

func GetDao(ctx *user.Ctx) *contentDao {
	if ctx == nil{
		log.Fatal("ctx can't nil")
	}
	return &contentDao{ctx}
}

func (d *contentDao) GetMomentListDB(db *gorm.DB,req *content.MomentListReq) ([]*content.Moment,error) {
	var moments []*content.Moment
	err := db.Where(`deleted_at = ?`, dbi.PostgreZeroTime).
		Limit(int(req.PageSize)).Offset(int((req.PageNo - 1) * req.PageSize)).
		Find(&moments).Error
	if err != nil {
		return nil, err
	}
	return moments,nil
}

func (d *contentDao) DeleteMomentDB(db *gorm.DB,id,userId uint64) error {
	return db.Table("moment").Where(`id = ? AND user_id = ? `, id,userId).
		UpdateColumns(dbi.DeleteAt(d.TimeString)).Error
}