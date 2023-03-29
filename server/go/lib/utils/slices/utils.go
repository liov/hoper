package slices

import (
	"github.com/liov/hoper/server/go/lib/utils/def/constraints"
)

func Contains[T comparable](arr []T, sub T) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == sub {
			return true
		}
	}
	return false
}

func In[T comparable](a T, b []T) bool {
	for _, x := range b {
		if x == a {
			return true
		}
	}
	return false
}

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

func SlicesToMap[T any, K comparable, V any](slices []T, getKV func(T) (K, V)) map[K]V {
	m := make(map[K]V)
	for _, s := range slices {
		k, v := getKV(s)
		m[k] = v
	}
	return m
}

func Swap[T any](heap []T, i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}
