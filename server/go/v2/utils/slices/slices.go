package slices

import (
	"errors"
	"reflect"
	"sort"

	"github.com/liov/hoper/go/v2/utils/def"
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
// 没有泛型，范例，实际需根据不同类型各写一遍,用CmpKey，基本类型又用不了，go需要能给基本类型实现方法不能给外部类型实现方法
func IsCoincide(s1, s2 []def.CmpKey) bool {
	for i := range s1 {
		for j := range s2 {
			if s1[i].CmpKey() ==  s2[j].CmpKey() {
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

type IntSlice struct {
	s	[]int
}

func (p *IntSlice) Len() int           { return len(p.s) }
func (p *IntSlice) Less(i, j int) bool { return p.s[i] < p.s[j] }
func (p *IntSlice) Swap(i, j int)      { p.s[i], p.s[j] = p.s[j], p.s[i] }

// Sort is a convenience method.
func (p *IntSlice) Sort() { sort.Sort(p) }

func Sort(s []int){
	slice := IntSlice{s: s}
	slice.Sort()
}