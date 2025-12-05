package db

import (
	"github.com/hopeio/gox/context/httpctx"
	"gorm.io/gorm"
)

type CommonDao struct {
	*httpctx.Context
	db *gorm.DB
}

func GetDao(ctx *httpctx.Context, d *gorm.DB) *CommonDao {
	return &CommonDao{
		Context: ctx,
		db:      d,
	}
}
