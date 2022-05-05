package gormi

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"gorm.io/gorm"
)

func DelDB(db *gorm.DB, tableName string, id uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = now()
WHERE id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `'`
	return db.Exec(sql, id).Error
}

func DelByAuthDB(db *gorm.DB, tableName string, id, userId uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = now()
WHERE id = ?  AND user_id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `'`
	return db.Exec(sql, id, userId).Error
}

func ExistsDB(db *gorm.DB, tableName string, id uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `' LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ExistsByAuthDB(db *gorm.DB, tableName string, id, userId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE id = ?  AND user_id = ? AND deleted_at = '` + dbi.PostgreZeroTime + `' LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, userId).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
