package dbi

const (
	PostgreZeroTime = "0001-01-01 00:00:00"
	MysqlZeroTime   = "0000-00-00 00:00:00"
	tmFmtWithMS     = "2006-01-02 15:04:05.999"
	tmFmtZero       = "0000-00-00 00:00:00"
	nullStr         = "NULL"
)

const (
	ColumnDeletedAt = "deleted_at"
	ColumnId        = "id"
)

const (
	PostgreNotDeleted = ColumnDeletedAt + " = '" + PostgreZeroTime + "'"
	MysqlNotDeleted   = ColumnDeletedAt + " = '" + MysqlZeroTime + "'"
)
