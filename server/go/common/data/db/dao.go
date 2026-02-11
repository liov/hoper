package db

import (
	"gorm.io/gorm"
)

type CommonDao struct {
	*gorm.DB
}

func GetDao(db *gorm.DB) *CommonDao {
	return &CommonDao{
		db,
	}
}
