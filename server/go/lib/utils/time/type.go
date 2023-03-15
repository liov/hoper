package timei

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	strings "github.com/liov/hoper/server/go/lib/utils/strings"
	"strconv"
	"time"
)

type Date time.Time

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	return date.Format(DateFormat), nil
}

func (date Date) Format(foramt string) string {
	return time.Time(date).Format(foramt)
}

// GormDataType gorm common data type
func (date *Date) GormDataType() string {
	return "date"
}

func (date Date) MarshalBinary() ([]byte, error) {
	return time.Time(date).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (date *Date) UnmarshalBinary(data []byte) error {
	return (*time.Time)(date).UnmarshalBinary(data)
}

func (date Date) GobEncode() ([]byte, error) {
	return date.MarshalBinary()
}

func (date *Date) GobDecode(data []byte) error {
	return date.UnmarshalBinary(data)
}

func (date Date) MarshalJSON() ([]byte, error) {
	t := time.Time(date)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(DateFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, DateFormat)
	b = append(b, '"')
	return b, nil
}

func (date *Date) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*date = (Date)(t)
	return err
}

type Time time.Time

func (dt *Time) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*dt = Time(nullTime.Time)
	return
}

func (dt Time) Value() (driver.Value, error) {
	return time.Time(dt), nil
}

func (dt Time) Format(foramt string) string {
	return time.Time(dt).Format(foramt)
}

// GormDataType gorm common data type
func (dt *Time) GormDataType() string {
	return "datetime"
}

func (dt Time) MarshalBinary() ([]byte, error) {
	return time.Time(dt).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (dt *Time) UnmarshalBinary(data []byte) error {
	return (*time.Time)(dt).UnmarshalBinary(data)
}

func (dt Time) GobEncode() ([]byte, error) {
	return dt.MarshalBinary()
}

func (dt *Time) GobDecode(data []byte) error {
	return dt.UnmarshalBinary(data)
}

func (dt Time) MarshalJSON() ([]byte, error) {
	t := time.Time(dt)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(TimeFormatDisplay)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, TimeFormatDisplay)
	b = append(b, '"')
	return b, nil
}

func (dt *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t, err := time.ParseInLocation(`"`+TimeFormatDisplay+`"`, string(data), time.Local)
	*dt = (Time)(t)
	return err
}

type UnixTime time.Time

func (ut *UnixTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ut = UnixTime(nullTime.Time)
	return
}

func (ut UnixTime) Value() (driver.Value, error) {
	return time.Time(ut), nil
}

func (ut UnixTime) Format(foramt string) string {
	return time.Time(ut).Format(foramt)
}

// GormDataType gorm common data type
func (ut *UnixTime) GormDataType() string {
	return "datetime"
}

func (ut UnixTime) MarshalBinary() ([]byte, error) {
	return time.Time(ut).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ut *UnixTime) UnmarshalBinary(data []byte) error {
	return (*time.Time)(ut).UnmarshalBinary(data)
}

func (ut UnixTime) GobEncode() ([]byte, error) {
	return ut.MarshalBinary()
}

func (ut *UnixTime) GobDecode(data []byte) error {
	return ut.UnmarshalBinary(data)
}

func (ut UnixTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ut)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	return strings.ToBytes(strconv.Itoa(int(t.Unix()))), nil
}

func (ut *UnixTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	str, err := strconv.Atoi(strings.ToString(data))
	if err != nil {
		return err
	}
	t := time.Unix(int64(str), 0)
	*ut = (UnixTime)(t)
	return err
}

type UnixNanoTime time.Time

func (unt *UnixNanoTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*unt = UnixNanoTime(nullTime.Time)
	return
}

func (unt UnixNanoTime) Value() (driver.Value, error) {
	return time.Time(unt), nil
}

func (unt UnixNanoTime) Format(foramt string) string {
	return time.Time(unt).Format(foramt)
}

// GormDataType gorm common data type
func (unt *UnixNanoTime) GormDataType() string {
	return "datetime"
}

func (unt UnixNanoTime) MarshalBinary() ([]byte, error) {
	return time.Time(unt).MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (unt *UnixNanoTime) UnmarshalBinary(data []byte) error {
	return (*time.Time)(unt).UnmarshalBinary(data)
}

func (unt UnixNanoTime) GobEncode() ([]byte, error) {
	return unt.MarshalBinary()
}

func (unt *UnixNanoTime) GobDecode(data []byte) error {
	return unt.UnmarshalBinary(data)
}

func (unt UnixNanoTime) MarshalJSON() ([]byte, error) {
	t := time.Time(unt)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	return strings.ToBytes(strconv.Itoa(int(t.UnixNano()))), nil
}

func (unt *UnixNanoTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	str, err := strconv.Atoi(strings.ToString(data))
	if err != nil {
		return err
	}
	t := time.Unix(0, int64(str))
	*unt = (UnixNanoTime)(t)
	return err
}

// 对应数据库datetime或timestamp,或date
// jsonType 0 序列化为 "2006-01-02 15:04:05",jsonType 1序列化为"2006-01-02",jsonType 2 序列化为秒时间戳, jsonType 3序列化为毫秒时间戳
// 序列化,反序列化前需设置typ
type UnionTime struct {
	time.Time
	jsonType uint8
}

func NewUnionTime(t time.Time, typ uint8) UnionTime {
	return UnionTime{Time: t, jsonType: typ}
}

func ZeroUnionTime(typ uint8) UnionTime {
	return UnionTime{jsonType: typ}
}
func (ut *UnionTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ut = UnionTime{Time: nullTime.Time}
	return
}

func (ut UnionTime) Value() (driver.Value, error) {
	if ut.jsonType == 1 {
		return ut.Format(DateFormat), nil
	}
	return ut.Time, nil
}

// GormDataType gorm common data type
func (ut *UnionTime) GormDataType() string {
	return "datetime"
}

func (ut UnionTime) Format(foramt string) string {
	return ut.Time.Format(foramt)
}

func (ut UnionTime) MarshalBinary() ([]byte, error) {
	return ut.Time.MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ut *UnionTime) UnmarshalBinary(data []byte) error {
	return ut.Time.UnmarshalBinary(data)
}

func (ut UnionTime) GobEncode() ([]byte, error) {
	return ut.MarshalBinary()
}

func (ut *UnionTime) GobDecode(data []byte) error {
	return ut.UnmarshalBinary(data)
}

func (ut UnionTime) MarshalJSON() ([]byte, error) {
	t := ut.Time
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	switch ut.jsonType {
	case 0:
		b := make([]byte, 0, len(TimeFormatDisplay)+2)
		b = append(b, '"')
		b = t.AppendFormat(b, TimeFormatDisplay)
		b = append(b, '"')
		return b, nil
	case 1:
		b := make([]byte, 0, len(DateFormat)+2)
		b = append(b, '"')
		b = t.AppendFormat(b, DateFormat)
		b = append(b, '"')
		return b, nil
	case 2:
		return strings.ToBytes(strconv.Itoa(int(t.Unix()))), nil
	case 3:
		return strings.ToBytes(strconv.Itoa(int(t.UnixNano()))), nil
	}

	return nil, nil
}

func (ut *UnionTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	var err error
	var t time.Time
	switch ut.jsonType {
	case 0:
		t, err = time.ParseInLocation(`"`+TimeFormatDisplay+`"`, string(data), time.Local)
	case 1:
		t, err = time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	case 2:
		str, err := strconv.Atoi(strings.ToString(data))
		if err != nil {
			return err
		}
		t = time.Unix(int64(str), 0)
	case 3:
		str, err := strconv.Atoi(strings.ToString(data))
		if err != nil {
			return err
		}
		t = time.Unix(0, int64(str))
	}
	*ut = UnionTime{Time: t, jsonType: ut.jsonType}
	return err
}

func (ut *UnionTime) Type(typ uint8) UnionTime {
	ut.jsonType = typ
	return *ut
}

type TimeStamp int64
type UnixTimeStamp int64

func (ts *UnixTimeStamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = UnixTimeStamp(nullTime.Time.Unix())
	return
}

func (ts UnixTimeStamp) Value() (driver.Value, error) {
	return time.Unix(int64(ts), 0), nil
}

func (ts UnixTimeStamp) Format(foramt string) string {
	return time.Unix(int64(ts), 0).Format(foramt)
}

// GormDataType gorm common data type
func (ts UnixTimeStamp) GormDataType() string {
	return "datetime"
}

type UnixNanoTimeStamp int64

func (ts *UnixNanoTimeStamp) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = UnixNanoTimeStamp(nullTime.Time.UnixNano())
	return
}

func (ts UnixNanoTimeStamp) Value() (driver.Value, error) {
	return time.Unix(0, int64(ts)), nil
}

func (ts UnixNanoTimeStamp) Format(foramt string) string {
	return time.Unix(0, int64(ts)).Format(foramt)
}

// GormDataType gorm common data type
func (ts UnixNanoTimeStamp) GormDataType() string {
	return "datetime"
}
