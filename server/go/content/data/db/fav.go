package db

import (
	sqlib "database/sql"

	"github.com/hopeio/scaffold/errcode"

	sqlx "github.com/hopeio/gox/database/sql"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *ContentDao) FavExists(title string) (uint64, error) {
	ctxi := d
	sql := `SELECT id FROM "` + model.TableNameFavorite + `" 
WHERE title = ? AND user_id = ?` + sqlx.WithNotDeleted
	var id uint64

	err := d.db.Raw(sql, title, d.AuthInfo).Scan(&id).Error
	if err != nil && err != sqlib.ErrNoRows {
		return 0, ctxi.RespErrorLog(errcode.DBError, err, "ContainerExists")
	}
	return id, nil
}
