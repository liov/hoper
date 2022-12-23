package postgres

import (
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db/const"
	"gorm.io/gorm"
)

const existsSQL = `SELECT EXISTS(SELECT * FROM %s WHERE %s = ?` + dbi.WithNotDeleted + ` LIMIT 1)`
const deleteSQL = `Update %s SET deleted_at = now() WHERE %s = ?` + dbi.WithNotDeleted

func Delete(db *gorm.DB, tableName string, id uint64) error {
	return DeleteSQL(db, tableName, "id", id)
}

func DeleteSQL(db *gorm.DB, tableName, column string, value any) error {
	sql := `Update ` + tableName + ` SET deleted_at = now()
WHERE ` + column + ` = ?` + dbi.WithNotDeleted
	return db.Exec(sql, value).Error
}

func DeleteByAuth(db *gorm.DB, tableName string, id, userId uint64) error {
	sql := `Update ` + tableName + ` SET deleted_at = now()
WHERE id = ?  AND user_id = ?` + dbi.WithNotDeleted
	return db.Exec(sql, id, userId).Error
}

func ExistsByIdWithDeletedAt(db *gorm.DB, tableName string, id uint64) (bool, error) {
	return ExistsBySQL(db, ExistsSQL(tableName, "id", true), id)
}

func ExistsByAuthWithDeletedAt(db *gorm.DB, tableName string, id, userId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM ` + tableName + ` 
WHERE id = ?  AND user_id = ?` + dbi.WithNotDeleted + ` LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, userId).Scan(&exists).Error
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
	sql := `SELECT EXISTS(SELECT * FROM ` + tableName + ` WHERE ` + column + ` = ?`
	if withDeletedAt {
		sql += dbi.WithNotDeleted
	}
	sql += ` LIMIT 1)`
	return sql
}

func ExistsBySQL(db *gorm.DB, sql string, value ...any) (bool, error) {
	var exists bool
	err := db.Raw(sql, value...).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func Exists(db *gorm.DB, tableName, column string, value interface{}, withDeletedAt bool) (bool, error) {
	return ExistsBySQL(db, ExistsSQL(tableName, column, withDeletedAt), value)
}
