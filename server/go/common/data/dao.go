package data

import (
	"github.com/liov/hoper/server/go/common/data/db"
	"gorm.io/gorm"
)

func GetDBDao(d *gorm.DB) *db.CommonDao {
	return db.GetDao(d)
}
