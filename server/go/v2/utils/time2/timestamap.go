package time2

import (
	"database/sql/driver"
	"strconv"
	"time"
)

//毫秒
type Time int64

// Scan scan time.
func (t *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*t = Time(sc.Unix()*1000 + sc.UnixNano()/1e6)
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*t = Time(i)
	}
	return
}

// Value get time value.
func (t Time) Value() (driver.Value, error) {
	s := t / 1000
	ns := (t % 1000) * 1e6
	return time.Unix(int64(s), int64(ns)), nil
}

// Time get time.
func (t Time) Time() time.Time {
	s := t / 1000
	ns := (t % 1000) * 1e6
	return time.Unix(int64(s), int64(ns))
}

func Time2(t time.Time) Time {
	return Time(t.UnixNano() / 1e6)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}
