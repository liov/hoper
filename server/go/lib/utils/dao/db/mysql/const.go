package mysql

import (
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db/const"
)

const (
	DateTimeZero  = "0001-01-01 00:00:00"
	TimeStampZero = "0000-00-00 00:00:00"
)

const (
	NotDeleted = dbi.ColumnDeletedAt + " = '" + DateTimeZero + "'"
)
