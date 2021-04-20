package slices

import (
	"errors"
	"reflect"
	"sort"
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

type Less interface {
	Less(interface{}) bool
}

type IntSlice struct {
	s []int
}

func (p *IntSlice) Len() int           { return len(p.s) }
func (p *IntSlice) Less(i, j int) bool { return p.s[i] < p.s[j] }
func (p *IntSlice) Swap(i, j int)      { p.s[i], p.s[j] = p.s[j], p.s[i] }

// Sort is a convenience method.
func (p *IntSlice) Sort() { sort.Sort(p) }

func Sort(s []int) {
	slice := IntSlice{s: s}
	slice.Sort()
}
