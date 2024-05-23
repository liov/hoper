package db

import (
	sqlib "database/sql"
	"github.com/hopeio/cherry/protobuf/errorcode"
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *ContentDao) FavExists(title string) (uint64, error) {
	ctxi := d
	sql := `SELECT id FROM "` + model.FavoritesTableName + `" 
WHERE title = ? AND user_id = ?` + dbi.WithNotDeleted
	var id uint64

	err := d.db.Raw(sql, title, d.AuthInfo).Scan(&id).Error
	if err != nil && err != sqlib.ErrNoRows {
		return 0, ctxi.ErrorLog(errorcode.DBError, err, "ContainerExists")
	}
	return id, nil
}
