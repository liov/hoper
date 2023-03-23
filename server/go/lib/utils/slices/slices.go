package slices

type Slices[T any] []T

func (slices Slices[T]) Len() int { return len(slices) }

func (slices Slices[T]) ForEach(fn func(T)) {
	for _, t := range slices {
		fn(t)
	}
}

func (slices Slices[T]) Filter(fn func(T) bool) []T {
	var newSlices []T
	for _, t := range slices {
		if fn(t) {
			newSlices = append(newSlices, t)
		}
	}
	return newSlices
}

func (slices Slices[T]) Every(fn func(T) bool) {
	for _, t := range slices {
		if !fn(t) {
			return
		}
	}
}

func (slices Slices[T]) Some(fn func(T) bool) {
	for _, t := range slices {
		if fn(t) {
			return
		}
	}
}

func (slices Slices[T]) Zip(s []T) [][2]T {
	var newSlices [][2]T
	for i := range slices {
		newSlices = append(newSlices, [2]T{slices[i], s[i]})
	}
	return newSlices
}

func (slices Slices[T]) Reduce(fn func(T, T) T) T {
	ret := fn(slices[0], slices[1])
	for i := 2; i < len(slices); i++ {
		ret = fn(ret, slices[i])
	}
	return ret
}

func Map[S ~[]T, T any, V any](slices S, fn func(T) V) []V {
	ret := make([]V, 0, len(slices))
	for _, s := range slices {
		ret = append(ret, fn(s))
	}
	return ret
}

func Reduce[S ~[]T, T any](slices S, fn func(T, T) T) T {
	ret := fn(slices[0], slices[1])
	for i := 2; i < len(slices); i++ {
		ret = fn(ret, slices[i])
	}
	return ret
}

// 学学kotlin的定义
type Array[S, T any] []S

//type Function[T any] func[T]()

func (a Array[S, T]) Map(fn func(S) T) []T {
	ret := make([]T, 0, len(a))
	for _, s := range a {
		ret = append(ret, fn(s))
	}
	return ret
}
