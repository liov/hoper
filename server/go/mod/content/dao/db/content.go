package db

import (
	sqlib "database/sql"
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"github.com/actliboy/hoper/server/go/mod/content/model"
)

func (d *ContentDBDao) FavExists(title string) (uint64, error) {
	ctxi := d
	sql := `SELECT id FROM "` + model.FavoritesTableName + `" 
WHERE title = ? AND user_id = ? AND ` + dbi.PostgreNotDeleted
	var id uint64

	err := d.db.Raw(sql, title, d.IdStr).Row().Scan(&id)
	if err != nil && err != sqlib.ErrNoRows {
		return 0, ctxi.ErrorLog(errorcode.DBError, err, "ContainerExists")
	}
	return id, nil
}
