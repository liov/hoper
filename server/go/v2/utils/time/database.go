package timei

import "time"

const (
	FormatTime = "2006-01-02 15:04:05.999999"
)

func PgNow() string {
	return time.Now().Format(time.RFC3339Nano)
}

func DBNow() string {
	return time.Now().Format(FormatTime)
}
