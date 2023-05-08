package db

import (
	"database/sql"
	"github.com/hopeio/pandora/context/http_context"
	"github.com/hopeio/pandora/protobuf/errorcode"
	clausei "github.com/hopeio/pandora/utils/dao/db/gorm/clause"
	"github.com/hopeio/pandora/utils/log"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContentDBDao struct {
	*http_context.Context
	ChainDao
}

func GetDao(ctx *http_context.Context, db *gorm.DB) *ContentDBDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDBDao{Context: ctx, ChainDao: ChainDao{db: db}}
}

func (d *ContentDBDao) SetDB(db *gorm.DB) {
	d.db = db
}

func (d *ContentDBDao) Begin() *ContentDBDao {
	return GetDao(d.Context, d.db.Begin())
}

type ChainDao struct {
	clausei.Clause2
	db *gorm.DB
}

func (d *ContentDBDao) CreateContextExt(typ content.ContentType, refId uint64) error {
	err := d.db.Exec(`INSERT INTO `+model.ContentExtTableName+`(type,ref_id) Values(?,?)`, typ, refId).Error
	if err != nil {
		return d.Context.ErrorLog(errorcode.DBError, err, "CreateContextExt")
	}
	return nil
}

func (d *ContentDBDao) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	err := d.db.Transaction(fc, opts...)
	if err != nil && err != errorcode.DBError {
		d.Context.Error(err.Error(), zap.String(log.Position, "Transaction"))
		return err
	}
	return err
}
