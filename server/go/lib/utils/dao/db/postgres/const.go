package postgres

import dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"

const (
	ZeroTime = "0001-01-01 00:00:00"
)

const (
	NotDeleted = dbi.ColumnDeletedAt + " = '" + ZeroTime + "'"
)

const (
	WithNotDeleted = ` AND ` + NotDeleted
)
