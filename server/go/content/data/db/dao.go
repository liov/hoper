package db

import (
	"database/sql"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/log"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContentDao struct {
	*httpctx.Context
	ChainDao
}

func GetDao(ctx *httpctx.Context, db *gorm.DB) *ContentDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDao{Context: ctx, ChainDao: ChainDao{db: db}}
}

func (d *ContentDao) SetDB(db *gorm.DB) {
	d.db = db
}

func (d *ContentDao) Begin() *ContentDao {
	return GetDao(d.Context, d.db.Begin())
}

func (d *ContentDao) CreateContextExt(typ content.ContentType, refId uint64) error {
	err := d.db.Exec(`INSERT INTO `+model.ContentExtTableName+`(type,ref_id) Values(?,?)`, typ, refId).Error
	if err != nil {
		return d.Context.ErrorLog(errorcode.DBError, err, "CreateContextExt")
	}
	return nil
}

func (d *ContentDao) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	err := d.db.Transaction(fc, opts...)
	if err != nil && err != errorcode.DBError {
		d.Context.Error(err.Error(), zap.String(log.FieldPosition, "Transaction"))
		return err
	}
	return err
}
