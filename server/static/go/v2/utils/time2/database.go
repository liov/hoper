package time2

import "time"

func PgNow() string {
	return time.Now().Format(time.RFC3339Nano)
}

func MsNow() string {
	return time.Now().Format("2006-01-02 15:04:05.999")
}
