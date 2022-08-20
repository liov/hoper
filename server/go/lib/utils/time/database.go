package timei

import "time"

func PgNow() string {
	return time.Now().Format(time.RFC3339Nano)
}

func DBNow() string {
	return time.Now().Format(FormatTime)
}
