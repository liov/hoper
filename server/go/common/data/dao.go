package data

import (
	"github.com/hopeio/context/httpctx"
	"github.com/liov/hoper/server/go/common/data/db"
	"gorm.io/gorm"
)

func GetDBDao(ctx *httpctx.Context, d *gorm.DB) *db.CommonDao {
	return db.GetDao(ctx, d)
}
