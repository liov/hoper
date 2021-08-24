package dao

import (
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	"github.com/liov/hoper/server/go/mod/content/model"
	"gorm.io/gorm"
)

func (d *contentDao) FavExists(db *gorm.DB, title string) (uint64, error) {
	ctxi := d
	sql := `SELECT id FROM "` + model.FavoritesTableName + `" 
WHERE title = ? AND user_id = ? AND ` + dbi.PostgreNotDeleted
	var id uint64
	err := db.Raw(sql, title, d.IdStr).Row().Scan(&id)
	if err != nil {
		return 0, ctxi.ErrorLog(errorcode.DBError, err, "ContainerExistsDB")
	}
	return id, nil
}
