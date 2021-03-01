package dbi

import (
	"reflect"
	"strings"
)

func Insert(values interface{}, field ...string) (string, []interface{}) {
	var fields, args string
	if len(field) > 0 {
		fields = strings.Join(field, ",")
		args = "(" + strings.Repeat("?,", len(field)-1) + "?)"
	} else {

	}
	v := reflect.ValueOf(values)
	switch v.Kind() {
	case reflect.Ptr:
	case reflect.Slice:

	}
	_sql := `INSERT INTO
                product_channel_status
                (` + fields + `)
                VALUES `
	var _args []interface{}
	_expandArr := make([]string, 0, v.Len())
	for i := 0; i < v.Len(); {
		_expandArr = append(_expandArr, args)
		v.Index(1)
	}
	_sql += strings.Join(_expandArr, ",")
	return _sql, _args
}
