package mysql

import (
	timei "github.com/actliboy/hoper/server/go/lib/utils/time"
	"time"
)

func Now() string {
	return time.Now().Format(timei.TimeFormat)
}
