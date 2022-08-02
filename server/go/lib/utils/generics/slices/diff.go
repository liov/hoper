package slices

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

func Diff[T comparable](a, b []T) []T {
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
