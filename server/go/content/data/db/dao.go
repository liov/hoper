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

func GetDao(db *gorm.DB) *ContentDao {
	return &ContentDao{db}
}

func (d *ContentDao) CreateContextExt(ctx context.Context, typ content.ContentType, refId uint64) error {
	err := d.Exec(`INSERT INTO `+model.TableNameStatistics+`(type,ref_id) Values(?,?)`, typ, refId).Error
	if err != nil {
		log.Error("CreateContextExt", zap.Error(err))
		return errcode.DBError.Wrap(err)
	}
	return nil
}
