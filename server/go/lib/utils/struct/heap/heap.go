package heap

import (
	_interface "github.com/liov/hoper/server/go/lib/utils/def/interface"
	"golang.org/x/exp/constraints"
)

type Heap[T _interface.CmpKey[V], V constraints.Ordered] []T

func (heap Heap[T, V]) Init() {
	// heapify
	n := len(heap)
	for i := n/2 - 1; i >= 0; i-- {
		heap.down(i, n)
	}
}

func (heap *Heap[T, V]) Push(x T) {
	h := *heap
	*heap = append(h, x)
	heap.up(len(h))
}

func (heap Heap[T, V]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && heap[j2].CmpKey() > heap[j1].CmpKey() {
			j = j2 // = 2*i + 2  // right child
		}
		if !(heap[j].CmpKey() > heap[i].CmpKey()) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		i = j
	}
	return i > i0
}

func (heap *Heap[T, V]) Pop() T {
	h := *heap
	n := len(h) - 1
	item := h[0]
	h[0], h[n] = h[n], h[0]
	h.down(0, n)
	*heap = h[:n]
	return item
}

func (heap *Heap[T, V]) First() T {
	return (*heap)[0]
}

func (heap *Heap[T, V]) Last() T {
	return (*heap)[len(*heap)-1]
}

func (heap *Heap[T, V]) Remove(i int) T {
	h := *heap
	n := len(h) - 1
	if n != i {
		h[i], h[n] = h[n], h[i]
		if !heap.down(i, n) {
			heap.up(i)
		}
	}
	return heap.Pop()
}

func (heap *Heap[T, V]) Fix(i int) {
	if !heap.down(i, len(*heap)) {
		heap.up(i)
	}
}

func (heap Heap[T, V]) up(j int) {

	for {
		i := (j - 1) / 2 // parent
		if i == j || !(heap[j].CmpKey() > heap[i].CmpKey()) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		j = i
	}
}
