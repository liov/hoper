package timei

import (
	"strconv"
	"strings"
	"time"
)

func WeiboParse(t string) time.Time {
	if t == "刚刚" {
		return time.Now()
	} else if strings.HasSuffix(t, "分钟前") {
		diff, _ := strconv.Atoi(t[:len(t)-len("分钟前")])
		return time.Now().Add(time.Duration(diff) * time.Minute)
	} else if strings.HasSuffix(t, "小时前") {
		diff, _ := strconv.Atoi(t[:len(t)-len("小时前")])
		return time.Now().Add(time.Duration(diff) * time.Minute)
	} else if strings.HasPrefix(t, "昨天") {
		hms := t[len("昨天"):]
		tt, _ := Parse(TimeFormatDisplay, time.Now().Format(DateFormat)+hms)
		return tt
	}
	tt, _ := Parse(DateFormat, t)
	return tt
}

func Parse(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}
