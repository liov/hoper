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

type Equal interface {
	IsEqual(interface{}) bool
}

func IsCoincide(s1, s2 []Equal) bool {
	for i := range s1 {
		for j := range s2 {
			if s1[i].IsEqual(s2[j]) {
				return true
			}
		}
	}
	return false
}

func RemoveDuplicates(s []int) []int {
	var m = make(map[int]struct{})
	for _, i := range s {
		m[i] = struct{}{}
	}
	s = s[:0]
	for k, _ := range m {
		s = append(s, k)
	}
	return s
}
