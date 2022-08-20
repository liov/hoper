package timei

import "time"

const (
	SecondsOfDay    = 24 * 60 * 60
	SecondsOfMinute = 60
	TimeDay         = SecondsOfDay * time.Second
)

const (
	FormatTime        = "2006-01-02 15:04:05.999999"
	DisplayFormatTime = "2006-01-02 15:04:05"
	SimpleFormatTime  = "2006-01-02 15:04:05"
	CompactFormatTime = "20060102150405"
)
