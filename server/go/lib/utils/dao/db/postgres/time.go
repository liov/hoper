package postgres

import "time"

func Now() string {
	return time.Now().Format(time.RFC3339Nano)
}
