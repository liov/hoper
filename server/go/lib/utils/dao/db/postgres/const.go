package postgres

import (
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db/const"
)

const (
	ZeroTimeUCT     = "0001-01-01 00:00:00"
	ZeroTimeUCTZone = ZeroTimeUCT + "+00:00:00"
	ZeroTimeCST     = "0001-01-01 08:05:43"
	ZeroTimeCSTZone = ZeroTimeCST + "+08:05:43"
)

const (
	NotDeletedUCT = dbi.ColumnDeletedAt + " = '" + ZeroTimeUCT + "'"
	NotDeletedCST = dbi.ColumnDeletedAt + " = '" + ZeroTimeCST + "'"
)

const (
	WithNotDeletedUCT = ` AND ` + NotDeletedUCT
	WithNotDeletedCST = ` AND ` + NotDeletedUCT
)
