package postgres

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db/const"
)

const (
	ZeroTime         = "0001-01-01 00:00:00"
	ZeroTimeTimeZone = "0001-01-01 08:05:43+08:05:43"
)

const (
	NotDeleted = dbi.ColumnDeletedAt + " = '" + ZeroTime + "'"
)

const (
	WithNotDeleted = ` AND ` + NotDeleted
)
