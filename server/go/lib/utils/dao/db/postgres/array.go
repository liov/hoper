package postgres

import (
	"bytes"
	"database/sql/driver"
	"strconv"
	"strings"
)

type IntArray []int

func (d *IntArray) Scan(value interface{}) error {
	str := value.(string)
	strs := strings.Split(str[1:len(str)-1], ",")
	var arr []int
	for _, numstr := range strs {
		num, err := strconv.Atoi(numstr)
		if err != nil {
			return err
		}
		arr = append(arr, num)
	}
	*d = arr
	return nil
}

func (d IntArray) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, num := range d {
		buf.WriteString(strconv.Itoa(num))
		if i != len(d)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}
