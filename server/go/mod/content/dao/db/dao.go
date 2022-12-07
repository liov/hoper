package db

import (
	"database/sql"
	contexti "github.com/liov/hoper/server/go/lib/context"
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	clausei "github.com/liov/hoper/server/go/lib/utils/dao/db/gorm/clause"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/mod/content/model"
	"github.com/liov/hoper/server/go/mod/protobuf/content"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContentDBDao struct {
	*contexti.Ctx
	ChainDao
}

func GetDao(ctx *contexti.Ctx, db *gorm.DB) *ContentDBDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &ContentDBDao{Ctx: ctx, ChainDao: ChainDao{db: db}}
}

func (d *ContentDBDao) SetDB(db *gorm.DB) {
	d.db = db
}

func (d *ContentDBDao) Begin() *ContentDBDao {
	return GetDao(d.Ctx, d.db.Begin())
}

type ChainDao struct {
	clausei.Clause2
	db *gorm.DB
}

func (d *ContentDBDao) CreateContextExt(typ content.ContentType, refId uint64) error {
	err := d.db.Exec(`INSERT INTO `+model.ContentExtTableName+`(type,ref_id) Values(?,?)`, typ, refId).Error
	if err != nil {
		return d.Ctx.ErrorLog(errorcode.DBError, err, "CreateContextExt")
	}
	return nil
}

func (d *ContentDBDao) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	err := d.db.Transaction(fc, opts...)
	if err != nil && err != errorcode.DBError {
		d.Ctx.Error(err.Error(), zap.String(log.Position, "Transaction"))
		return err
	}
	return err
}
