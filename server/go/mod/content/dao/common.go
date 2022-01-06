package dao

import (
	"database/sql"
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/mod/content/model"
	"github.com/actliboy/hoper/server/go/mod/protobuf/content"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (d *contentDao) CreateContextExtDB(db *gorm.DB, typ content.ContentType, refId uint64) error {
	err := db.Exec(`INSERT INTO `+model.ContentExtTableName+`(type,ref_id) Values(?,?)`, typ, refId).Error
	if err != nil {
		return d.Ctx.ErrorLog(errorcode.DBError, err, "CreateContextExtDB")
	}
	return nil
}

func (d *contentDao) Transaction(db *gorm.DB, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	err := db.Transaction(fc, opts...)
	if err != nil && err != errorcode.DBError {
		d.Ctx.Error(err.Error(), zap.String(log.Position, "Transaction"))
		return err
	}
	return err
}
