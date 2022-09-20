package timei

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// 毫秒
type TimePs int64

// Scan scan time.
func (t *TimePs) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*t = TimePs(sc.Unix()*1000 + sc.UnixNano()/1e6)
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*t = TimePs(i)
	}
	return
}

// Value get time value.
func (t TimePs) Value() (driver.Value, error) {
	s := t / 1000
	ns := (t % 1000) * 1e6
	return time.Unix(int64(s), int64(ns)), nil
}

// Time get time.
func (t TimePs) Time() time.Time {
	s := t / 1000
	ns := (t % 1000) * 1e6
	return time.Unix(int64(s), int64(ns))
}

func Time2(t time.Time) TimePs {
	return TimePs(t.UnixNano() / 1e6)
}

func (t TimePs) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}
