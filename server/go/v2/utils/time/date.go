package timei

import (
	"strconv"
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
func GetYMD(time time.Time,sep string) string {
	year, month, day := time.Date()

	var monthStr string
	var dateStr string
	if month < 10 {
		monthStr = "0" + strconv.Itoa(int(month+1))
	} else {
		monthStr = strconv.Itoa(int(month + 1))
	}

	if day < 10 {
		dateStr = "0" + strconv.Itoa(day)
	} else {
		dateStr = strconv.Itoa(day)
	}
	return strconv.Itoa(year) + sep + monthStr + sep + dateStr
}

// GetYM 得到以sep为分隔符的年、月字符串(今天所属于的月份)
func GetYM(time time.Time,sep string) string {
	year, month, _ := time.Date()

	var monthStr string
	if month < 10 {
		monthStr = "0" + strconv.Itoa(int(month+1))
	} else {
		monthStr = strconv.Itoa(int(month + 1))
	}
	return strconv.Itoa(year) + sep + monthStr
}

// GetYesterdayYMD 得到以sep为分隔符的年、月、日字符串(昨天)
func GetYesterdayYMD(sep string) string {
	return GetYM(time.Now().AddDate(0,0,-1),sep)
}

// GetTomorrowYMD 得到以sep为分隔符的年、月、日字符串(明天)
func GetTomorrowYMD(sep string) string {
	return GetYM(time.Now().AddDate(0,0,1),sep)
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
	return GetTodayZeroTime().AddDate(0,0,-1)
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
