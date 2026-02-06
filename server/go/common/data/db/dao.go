package db

import (
	"github.com/hopeio/gox/context/httpctx"
	"gorm.io/gorm"
)

type CommonDao struct {
	context.Context
	db *gorm.DB
}

func GetDao(ctx context.Context, d *gorm.DB) *CommonDao {
	return &CommonDao{
		Context: ctx,
		db:      d,
	}
}
