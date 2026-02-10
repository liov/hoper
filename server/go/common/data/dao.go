package data

import (
	"context"

	"github.com/liov/hoper/server/go/common/data/db"
	"gorm.io/gorm"
)

func GetDBDao(ctx context.Context, d *gorm.DB) *db.CommonDao {
	return db.GetDao(ctx, d)
}
