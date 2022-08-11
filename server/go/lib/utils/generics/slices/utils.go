package slices

import "github.com/actliboy/hoper/server/go/lib/utils/generics/def/constraints"

func ReverseRunes[T any](runes []T) []T {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return runes
}

func Max[T constraints.Number](s []T) T {
	if len(s) == 0 {
		return *new(T)
	}
	max := s[0]
	if len(s) == 1 {
		return max
	}
	for i := 1; i < len(s); i++ {
		if s[i] > max {
			max = s[i]
		}
	}

	return max
}

func Min[T constraints.Number](s []T) T {
	if len(s) == 0 {
		return *new(T)
	}
	min := s[0]
	if len(s) == 1 {
		return min
	}
	for i := 1; i < len(s); i++ {
		if s[i] < min {
			min = s[i]
		}
	}

	return min
}
