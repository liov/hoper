package data

import (
	"gorm.io/gorm"
)

type uploadDao struct {
	*gorm.DB
}

func GetDao(db *gorm.DB) uploadDao {
	return uploadDao{db}
}
