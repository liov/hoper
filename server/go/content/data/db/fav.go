package db

import (
	sqlib "database/sql"

	"github.com/hopeio/scaffold/errcode"

	sqlx "github.com/hopeio/gox/database/sql"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *ContentDao) FavExists(title string, userId uint64) (uint64, error) {
	sql := `SELECT id FROM "` + model.TableNameFavorite + `" 
WHERE title = ? AND user_id = ?` + sqlx.WithNotDeleted
	var id uint64

	err := d.Raw(sql, title, userId).Scan(&id).Error
	if err != nil && err != sqlib.ErrNoRows {
		return 0, errcode.DBError.Wrap(err)
	}
	return id, nil
}
