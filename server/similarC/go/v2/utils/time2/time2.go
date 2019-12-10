package time2

import (
	"context"
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
	return Time(t.UnixNano()/1e6)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return t.Time().MarshalJSON()
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < time.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(d))
	return d, ctx, cancel
}
