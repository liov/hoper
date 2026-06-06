package data

import (
	"gorm.io/gorm"
)

type chatDao struct {
	*gorm.DB
}

func GetDao(db *gorm.DB) *chatDao {
	return &chatDao{db}
}
