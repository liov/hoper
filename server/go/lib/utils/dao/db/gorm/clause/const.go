package clausei

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
)

const (
	IDEQUAL   = dbi.ColumnId + dbi.ExprEqual
	NAMEEQUAL = dbi.ColumnName + dbi.ExprEqual
)
