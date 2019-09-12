package slices

import (
	"errors"
	"reflect"
)

func Contains(arr interface{}, sub interface{}) (bool, error) {
	valueOf := reflect.ValueOf(arr)
	if valueOf.Kind() != reflect.Slice {
		return false, errors.New("错误的参数，第一个参数必须为切片类型")
	}
	subValue := reflect.ValueOf(sub)
	for i := 0; i < valueOf.Len(); i++ {
		if valueOf.Index(i) == subValue {
			return true, nil
		}
	}
	return false, nil
}

func StringContains(arr []string, sub string) bool {
	for _, v := range arr {
		if v == sub {
			return true
		}
	}
	return false
}
