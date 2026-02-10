package db

import (
	"context"

	"github.com/hopeio/scaffold/errcode"

	"github.com/hopeio/gox/log"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContentDao struct {
	*gorm.DB
}

func GetDao(ctx context.Context, db *gorm.DB) *ContentDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDao{db.Session(&gorm.Session{Context: ctx, NewDB: true})}
}

func (d *ContentDao) CreateContextExt(typ content.ContentType, refId uint64) error {
	err := d.Exec(`INSERT INTO `+model.TableNameStatistics+`(type,ref_id) Values(?,?)`, typ, refId).Error
	if err != nil {
		log.Error("CreateContextExt", zap.Error(err))
		return errcode.DBError.Wrap(err)
	}
	return nil
}
