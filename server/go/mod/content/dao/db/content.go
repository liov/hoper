package db

import (
	sqlib "database/sql"
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db/const"
	"github.com/liov/hoper/server/go/mod/content/model"
)

func (d *ContentDBDao) FavExists(title string) (uint64, error) {
	ctxi := d
	sql := `SELECT id FROM "` + model.FavoritesTableName + `" 
WHERE title = ? AND user_id = ?` + dbi.WithNotDeleted
	var id uint64

	err := d.db.Raw(sql, title, d.ID).Scan(&id).Error
	if err != nil && err != sqlib.ErrNoRows {
		return 0, ctxi.ErrorLog(errorcode.DBError, err, "ContainerExists")
	}
	return id, nil
}
