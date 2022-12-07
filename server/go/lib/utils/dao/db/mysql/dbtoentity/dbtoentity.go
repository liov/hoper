package dbtoentity

import (
	"database/sql"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	"github.com/liov/hoper/server/go/lib/utils/dao/db/mysql"
)

func MysqlConvert(db *sql.DB, filename string) {
	mysqlgen := mysqlgen{db: db}
	dbi.Convert(&mysqlgen, filename)
}

func MysqlConvertByTable(db *sql.DB, tableName string) {
	mysqlgen := mysqlgen{db: db}
	dbi.ConvertByTable(&mysqlgen, tableName)
}

type mysqlgen struct {
	db *sql.DB
}

func (m *mysqlgen) Tables() []string {
	var tables []string
	rows, _ := m.db.Query(`SHOW TABLES`)
	for rows.Next() {
		var table string
		rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables
}

func (m *mysqlgen) Fields(tableName string) []*dbi.Field {
	var dbfields []*dbi.Field
	rows, _ := m.db.Query(`SHOW FULL COLUMNS FROM ` + tableName)
	for rows.Next() {
		var dbfield dbi.Field
		rows.Scan(&dbfield.Field, &dbfield.Type, &dbfield.Comment)
		dbfields = append(dbfields, &dbfield)
	}
	return dbfields
}

func (m *mysqlgen) TypeToGoTYpe(typ string) string {
	return mysql.MySqlTypeToGoType(typ)
}
