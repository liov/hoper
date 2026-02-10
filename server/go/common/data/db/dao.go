package db

import (
	"context"

	"gorm.io/gorm"
)

type CommonDao struct {
	*gorm.DB
}

func GetDao(ctx context.Context, db *gorm.DB) *CommonDao {
	return &CommonDao{
		db.Session(&gorm.Session{Context: ctx, NewDB: true}),
	}
}
