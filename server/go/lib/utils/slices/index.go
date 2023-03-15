package slices

import "sync"

type Index[T any, O comparable] struct {
	idx   []O
	value []T
	sync.RWMutex
}

func NewIndex[T any, O comparable]() *Index[T, O] {
	return &Index[T, O]{
		idx:   make([]O, 0),
		value: make([]T, 0),
	}
}

func (i *Index[T, O]) Add(idx O, res T) {
	i.idx = append(i.idx, idx)
	i.value = append(i.value, res)
}

func (i *Index[T, O]) Get(idx O) T {
	i.RLock()
	defer i.RUnlock()
	for j, v := range i.idx {
		if v == idx {
			return i.value[j]
		}
	}
	return *new(T)
}

func (i *Index[T, O]) Set(idx O, v T) {
	i.RLock()
	defer i.RUnlock()
	for j, x := range i.idx {
		if x == idx {
			i.value[j] = v
		}
		return
	}
}

func (i *Index[T, O]) Remove(idx O) {
	i.Lock()
	defer i.Unlock()
	for j, v := range i.idx {
		if v == idx {
			i.idx = append(i.idx[0:j], i.idx[j:]...)
			i.value = append(i.value[0:j], i.value[j:]...)
			return
		}
	}
	return
}
