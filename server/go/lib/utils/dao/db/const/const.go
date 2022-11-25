package _const

const (
	TmFmtWithMS = "2006-01-02 15:04:05.999"
	TmFmtZero   = "0000-00-00 00:00:00"
	NullStr     = "NULL"
)

const (
	ColumnDeletedAt = "deleted_at"
	ColumnId        = "id"
	ColumnName      = "name"
)

const (
	ExprEqual    = " = ?"
	ExprNotEqual = " != ?"
	ExprGreater  = " > ?"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	SQLite   = "sqlite3"
)
