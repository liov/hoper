package db

import (
	"database/sql"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/scaffold/errcode"

	"github.com/hopeio/utils/log"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContentDao struct {
	*httpctx.Context
	db *gorm.DB
}

func GetDao(ctx *httpctx.Context, db *gorm.DB) *ContentDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDao{Context: ctx, db: db}
}

func (d *ContentDao) SetDB(db *gorm.DB) {
	d.db = db
}

func (d *ContentDao) CreateContextExt(typ content.ContentType, refId uint64) error {
	err := d.db.Exec(`INSERT INTO `+model.TableNameStatistics+`(type,ref_id) Values(?,?)`, typ, refId).Error
	if err != nil {
		return d.Context.RespErrorLog(errcode.DBError, err, "CreateContextExt")
	}
	return nil
}

func (d *ContentDao) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	err := d.db.Transaction(fc, opts...)
	if err != nil && err != errcode.DBError {
		d.Context.ErrorLog(err, zap.String(log.FieldPosition, "Transaction"))
		return err
	}
	return err
}
