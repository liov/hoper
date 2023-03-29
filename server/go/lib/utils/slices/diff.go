package slices

import "golang.org/x/exp/constraints"

// 没有泛型，范例，实际需根据不同类型各写一遍,用CmpKey，基本类型又用不了，go需要能给基本类型实现方法不能给外部类型实现方法
func IsCoincide[T comparable](s1, s2 []T) bool {
	for i := range s1 {
		for j := range s2 {
			if s1[i] == s2[j] {
				return true
			}
		}
	}
	return false
}

func RemoveDuplicates[T comparable](s []T) []T {
	var m = make(map[T]struct{})
	for _, i := range s {
		m[i] = struct{}{}
	}
	s = s[:0]
	for k, _ := range m {
		s = append(s, k)
	}
	return s
}

// 取并集
func Intersection[T comparable](a []T, b []T) []T {
	if len(a) < 64 && len(b) < 64 {
		if len(a) > len(b) {
			return intersection(a, b)
		}
		return intersection(b, a)
	}
	panic("TODO:大数组利用map取并集")
}

func intersection[T comparable](a []T, b []T) []T {
	var ret []T
	for _, x := range a {
		if In(x, b) {
			ret = append(ret, x)
		}
	}
	return ret
}

// 有序数组取交集
func OrderedArrayIntersection[T constraints.Ordered](a []T, b []T) []T {
	var ret []T
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	var idx int
	for _, x := range a {
		if x > b[len(b)-1] {
			return ret
		}
		for j := idx; idx < len(b)-1; j++ {
			if a[len(a)-1] < b[idx] {
				return ret
			}
			if x == b[idx] {
				ret = append(ret, x)
				idx = j
			}
		}
	}
	return ret
}

// 取并集
func Union[T comparable](a []T, b []T) []T {
	var m = make(map[T]struct{}, len(a)+len(b))
	for _, x := range a {
		m[x] = struct{}{}
	}
	for _, x := range b {
		m[x] = struct{}{}
	}
	var ret = make([]T, len(m))
	for k, _ := range m {
		ret = append(ret, k)
	}
	return ret
}

// 取差集
func Difference[T comparable](a []T, b []T) []T {
	if len(a) < 64 && len(b) < 64 {
		if len(a) > len(b) {
			return difference(a, b)
		}
		return difference(b, a)
	}
	panic("TODO:大数组利用map取差集")
}

func difference[T comparable](a []T, b []T) []T {
	var ret []T
	for _, x := range a {
		if !In(x, b) {
			ret = append(ret, x)
		}
	}
	return ret
}

// 取差集
func Difference2[T comparable](a, b []T) []T {
	if len(a) > len(b) {
		a, b = b, a
	}
	var diff []T
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
