package dao

import (
	"github.com/liov/hoper/v2/content/model"
	"github.com/liov/hoper/v2/protobuf/utils/errorcode"
	dbi "github.com/liov/hoper/v2/utils/dao/db"
	"gorm.io/gorm"
)

func (d *contentDao) FavExists(db *gorm.DB, title string) (uint64, error) {
	ctxi := d.ctxi
	sql := `SELECT id FROM "` + model.FavoritesTableName + `" 
WHERE title = ? AND user_id = ? AND ` + dbi.PostgreNotDeleted
	var id uint64
	err := db.Raw(sql, title, d.ctxi.IdStr).Row().Scan(&id)
	if err != nil {
		return 0, ctxi.ErrorLog(errorcode.DBError, err, "ContainerExistsDB")
	}
	return id, nil
}
