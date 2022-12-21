package timei

import "time"

const (
	SecondsOfDay    = 24 * 60 * 60
	SecondsOfMinute = 60
	TimeDay         = SecondsOfDay * time.Second
)

const (
	TimeFormat            = "2006-01-02 15:04:05.999999"
	TimeFormatDisplay     = "2006-01-02 15:04:05"
	TimeFormatPostgresDB  = time.RFC3339
	TimeFormatNoDate      = "15:04:05"
	DayEndTime            = "23:59:59"
	DayEndTimeWithSpace   = " 23:59:59"
	DayBeginTime          = "00:00:00"
	DayBeginTimeWithSpace = " 00:00:00"
	DateFormat            = "2006-01-02"
	TimeFormatCompact     = "20060102150405"
	TimeFormatRFC1        = "2006/01/02 - 15:04:05"
)
