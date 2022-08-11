package mysql

import dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"

const (
	ZeroTime = "0000-00-00 00:00:00"
)

const (
	NotDeleted = dbi.ColumnDeletedAt + " = '" + ZeroTime + "'"
)
