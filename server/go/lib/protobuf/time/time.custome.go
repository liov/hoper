package time

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

func (ts *Time) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*ts = Time{T: nullTime.Time.UnixNano()}
	return
}

func (ts Time) Value() (driver.Value, error) {
	return time.Unix(0, ts.T), nil
}

func (ts Time) Format(foramt string) string {
	return time.Unix(0, ts.T).Format(foramt)
}

// GormDataType gorm common data type
func (ts Time) GormDataType() string {
	return "datetime"
}

func (ts Time) MarshalBinary() ([]byte, error) {
	enc := []byte{
		byte(ts.T >> 56), // bytes 1-8: seconds
		byte(ts.T >> 48),
		byte(ts.T >> 40),
		byte(ts.T >> 32),
		byte(ts.T >> 24),
		byte(ts.T >> 16),
		byte(ts.T >> 8),
		byte(ts.T),
	}
	return enc, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *Time) UnmarshalBinary(data []byte) error {
	ts.T = int64(data[7]) | int64(data[6])<<8 | int64(data[5])<<16 | int64(data[4])<<24 |
		int64(data[3])<<32 | int64(data[2])<<40 | int64(data[1])<<48 | int64(data[0])<<56
	return nil
}

func (date Time) GobEncode() ([]byte, error) {
	return date.MarshalBinary()
}

func (date *Time) GobDecode(data []byte) error {
	return date.UnmarshalBinary(data)
}

func (date Time) MarshalJSON() ([]byte, error) {
	t := time.Unix(0, date.T)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.DateTime)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.DateTime)
	b = append(b, '"')
	return b, nil
}

func (date *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t, err := time.ParseInLocation(`"2006-01-02"`, string(data), time.Local)
	*date = Time{T: t.UnixNano()}
	return err
}
