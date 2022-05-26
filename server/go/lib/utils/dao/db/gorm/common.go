package gormi

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"gorm.io/gorm"
)

const WithPostgreNotDeleted = ` AND ` + dbi.PostgreNotDeleted

func Delete(db *gorm.DB, tableName string, id uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = now()
WHERE id = ?` + WithPostgreNotDeleted
	return db.Exec(sql, id).Error
}

func DeleteByAuth(db *gorm.DB, tableName string, id, userId uint64) error {
	sql := `Update "` + tableName + `" SET deleted_at = now()
WHERE id = ?  AND user_id = ?` + WithPostgreNotDeleted
	return db.Exec(sql, id, userId).Error
}

func ExistsByIdWithDeletedAt(db *gorm.DB, tableName string, id uint64) (bool, error) {
	return ExistsBySQL(db, ExistsSQL(tableName, "id", true), id)
}

func ExistsByAuthWithDeletedAt(db *gorm.DB, tableName string, id, userId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" 
WHERE id = ?  AND user_id = ?` + WithPostgreNotDeleted + ` LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, userId).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ExistsById(db *gorm.DB, tableName string, id uint64) (bool, error) {
	return ExistsBySQL(db, ExistsSQL(tableName, "id", false), id)
}

func ExistsByColumn(db *gorm.DB, tableName, column string, value any) (bool, error) {
	return ExistsBySQL(db, ExistsSQL(tableName, column, false), value)
}

func ExistsSQL(tableName, column string, withDeletedAt bool) string {
	sql := `SELECT EXISTS(SELECT * FROM "` + tableName + `" WHERE ` + column + ` = ?`
	if withDeletedAt {
		sql += WithPostgreNotDeleted
	}
	sql += ` LIMIT 1)`
	return sql
}

func ExistsBySQL(db *gorm.DB, sql string, value any) (bool, error) {
	var exists bool
	err := db.Raw(sql, value).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
