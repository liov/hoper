package timei

import (
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006/01/02 - 15:04:05"

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
func GetTodayYMD(sep string) string {
	now := time.Now()
	year := now.Year()
	month := StrToIntMonth(now.Month().String())
	date := now.Day()

	var monthStr string
	var dateStr string
	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}

	if date < 10 {
		dateStr = "0" + strconv.Itoa(date)
	} else {
		dateStr = strconv.Itoa(date)
	}
	return strconv.Itoa(year) + sep + monthStr + sep + dateStr + sep
}

// GetTodayYM 得到以sep为分隔符的年、月字符串(今天所属于的月份)
func GetTodayYM(sep string) string {
	now := time.Now()
	year := now.Year()
	month := StrToIntMonth(now.Month().String())

	var monthStr string
	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}
	return strconv.Itoa(year) + sep + monthStr + sep
}

// GetYesterdayYMD 得到以sep为分隔符的年、月、日字符串(昨天)
func GetYesterdayYMD(sep string) string {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	todaySec := today.Unix()            //秒
	yesterdaySec := todaySec - 24*60*60 //秒
	yesterdayTime := time.Unix(yesterdaySec, 0)
	yesterdayYMD := yesterdayTime.Format("2006-01-02")
	return strings.Replace(yesterdayYMD, "-", sep, -1)
}

// GetTomorrowYMD 得到以sep为分隔符的年、月、日字符串(明天)
func GetTomorrowYMD(sep string) string {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	todaySec := today.Unix()           //秒
	tomorrowSec := todaySec + 24*60*60 //秒
	tomorrowTime := time.Unix(tomorrowSec, 0)
	tomorrowYMD := tomorrowTime.Format("2006-01-02")
	return strings.Replace(tomorrowYMD, "-", sep, -1)
}

// GetTodayTime 返回今天零点的time
func GetTodayTime() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return today
}

// GetYesterdayTime 返回昨天零点的time
func GetYesterdayTime() time.Time {
	now := time.Now()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	yesterdaySec := today.Unix() - 24*60*60
	return time.Unix(yesterdaySec, 0)
}

type DateFilter struct {
	DateStart string `json:"dateStart" explain:"起始时间"`
	DateEnd   string `json:"dateEnd" explain:"结束时间"`
	RangeEnum int    `json:"rangeEnum" explain:"1-今天,2-本周，3-本月，4-今年"`
}

//赋值本周期，并返回下周期日期
func (d *DateFilter) Scope() (string, string) {
	lastStartTime, lastEndTime := d.scope()
	start := lastStartTime.Format("2006-01-02")
	end := lastEndTime.Format("2006-01-02")
	return start, end
}

func (d *DateFilter) scope() (time.Time, time.Time) {
	if d.DateStart != "" && d.DateEnd != "" {
		start, _ := time.Parse("2006-01-02", d.DateStart)
		end, _ := time.Parse("2006-01-02", d.DateEnd)
		diff := end.Sub(start) / (24 * time.Hour)
		return start.AddDate(0, 0, -int(diff)), start
	}
	//如果传的是RangeEnum，截止日期都是这一天
	now := time.Now()
	endTime := now.AddDate(0, 0, 1)
	end := endTime.Format("2006-01-02")
	d.DateEnd = end
	switch d.RangeEnum {
	case 1:
		start := now.Format("2006-01-02")
		d.DateStart = start
		return now.AddDate(0, 0, -1), endTime.AddDate(0, 0, -1)
	case 2:
		var weekday int
		weekday = int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		} else {
			weekday += 1
		}
		startTime := endTime.AddDate(0, 0, -weekday)
		d.DateStart = startTime.Format("2006-01-02")
		return startTime.AddDate(0, 0, -7), startTime
	case 3:
		day := now.Day()
		startTime := endTime.AddDate(0, 0, -day)
		d.DateStart = startTime.Format("2006-01-02")
		return startTime.AddDate(0, -1, 0), startTime
	case 4:
		year := strconv.Itoa(now.Year())
		d.DateStart = year + "-01-01"
		startTime, _ := time.Parse("2006-01-02", d.DateStart)
		return startTime.AddDate(-1, 0, 0), startTime
	}
	return now, now
}
