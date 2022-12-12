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
		tt, _ := time.Parse(TimeFormatDisplay, time.Now().Format(DateFormat)+hms)
		return tt
	}
	tt, _ := time.Parse(DateFormat, t)
	return tt
}
