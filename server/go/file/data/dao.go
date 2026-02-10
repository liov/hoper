package data

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type uploadDao struct {
	*gorm.DB
}

func GetDao(ctx context.Context, db *gorm.DB) uploadDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return uploadDao{db.Session(&gorm.Session{Context: ctx, NewDB: true})}
}
