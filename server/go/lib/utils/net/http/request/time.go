package request

import (
	timei "github.com/actliboy/hoper/server/go/lib/utils/time"
	"time"
)

type RequestAt struct {
	Time       time.Time
	TimeStamp  int64
	TimeString string
}

func (r *RequestAt) String() string {
	return r.TimeString
}

func NewRequestAt() *RequestAt {
	now := time.Now()
	return &RequestAt{
		Time:       now,
		TimeStamp:  now.Unix(),
		TimeString: now.Format(timei.TimeFormat),
	}
}

func NewRequestAtByTime(t time.Time) *RequestAt {
	return &RequestAt{
		Time:       t,
		TimeStamp:  t.Unix(),
		TimeString: t.Format(timei.TimeFormat),
	}
}
