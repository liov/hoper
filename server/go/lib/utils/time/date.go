package timei

import (
	"strconv"
	"time"
)

// StrToIntMonth 字符串月份转整数月份
func StrToIntMonth(month string) int {
	var data = map[string]int{
		"January":   0,
		"February":  1,
		"March":     2,
		"April":     3,
		"May":       4,
		"June":      5,
		"July":      6,
		"August":    7,
		"September": 8,
		"October":   9,
		"November":  10,
		"December":  11,
	}
	return data[month]
}

// GetTodayYMD 得到以sep为分隔符的年、月、日字符串(今天)
func GetYMD(time time.Time, sep string) string {
	year, month, day := time.Date()

	var monthStr string
	var dateStr string
	if month < 10 {
		monthStr = "0" + strconv.Itoa(int(month))
	} else {
		monthStr = strconv.Itoa(int(month))
	}

	if day < 10 {
		dateStr = "0" + strconv.Itoa(day)
	} else {
		dateStr = strconv.Itoa(day)
	}
	return strconv.Itoa(year) + sep + monthStr + sep + dateStr
}

// GetYM 得到以sep为分隔符的年、月字符串(今天所属于的月份)
func GetYM(time time.Time, sep string) string {
	year, month, _ := time.Date()

	var monthStr string
	if month < 10 {
		monthStr = "0" + strconv.Itoa(int(month))
	} else {
		monthStr = strconv.Itoa(int(month))
	}
	return strconv.Itoa(year) + sep + monthStr
}

// GetYesterdayYMD 得到以sep为分隔符的年、月、日字符串(昨天)
func GetYesterdayYMD(sep string) string {
	return GetYM(time.Now().AddDate(0, 0, -1), sep)
}

// GetTomorrowYMD 得到以sep为分隔符的年、月、日字符串(明天)
func GetTomorrowYMD(sep string) string {
	return GetYM(time.Now().AddDate(0, 0, 1), sep)
}

// GetTodayZeroTime 返回今天零点的time
func GetTodayZeroTime() time.Time {
	year, month, day := time.Now().Date()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return today
}

// GetYesterdayZeroTime 返回昨天零点的time
func GetYesterdayZeroTime() time.Time {
	return GetTodayZeroTime().AddDate(0, 0, -1)
}
