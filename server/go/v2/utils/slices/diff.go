package slices

import "github.com/liov/hoper/go/v2/utils/def"

// 没有泛型，范例，实际需根据不同类型各写一遍,用CmpKey，基本类型又用不了，go需要能给基本类型实现方法不能给外部类型实现方法
func IsCoincide(s1, s2 []def.CmpKey) bool {
	for i := range s1 {
		for j := range s2 {
			if s1[i].CmpKey() == s2[j].CmpKey() {
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

func DiffUint64(a, b []uint64) []uint64 {
	var diff []uint64
Loop:
	for _, i := range b {
		for _, j := range a {
			if i == j {
				continue Loop
			}
		}
		diff = append(diff, i)
	}
	return diff
}
