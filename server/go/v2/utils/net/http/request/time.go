package request

import "time"

type RequestAt struct {
	Time       time.Time
	TimeStamp  int64
	TimeString string
}
