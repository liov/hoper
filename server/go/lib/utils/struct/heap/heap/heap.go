package heap

import (
	"github.com/liov/hoper/server/go/lib/utils/struct/heap"
)

type Interface interface {
	CmpKey() int
}
type Heap[T Interface] heap.Heap[T, int]

func New[T Interface](l int) Heap[T] {
	maxHeap := make(Heap[T], 0, l)
	return maxHeap
}

func NewFromArray[T Interface](arr []T) Heap[T] {
	heap := Heap[T](arr)
	for i := 1; i < len(arr); i++ {
		heap.up(i)
	}
	return arr
}

func (heap Heap[T]) Init() {
	// heapify
	n := len(heap)
	for i := n/2 - 1; i >= 0; i-- {
		heap.down(i, n)
	}
}

func (heap *Heap[T]) Push(x T) {
	h := *heap
	*heap = append(h, x)
	heap.up(len(h))
}

func (heap Heap[T]) down(i0, n int) bool {
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

func (heap *Heap[T]) Pop() T {
	h := *heap
	n := len(h) - 1
	item := h[0]
	h[0], h[n] = h[n], h[0]
	h.down(0, n)
	*heap = h[:n]
	return item
}

func (heap *Heap[T]) First() T {
	return (*heap)[0]
}

func (heap *Heap[T]) Last() T {
	return (*heap)[len(*heap)-1]
}

func (heap *Heap[T]) Remove(i int) T {
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

func (heap *Heap[T]) Fix(i int) {
	if !heap.down(i, len(*heap)) {
		heap.up(i)
	}
}

func (heap Heap[T]) up(j int) {

	for {
		i := (j - 1) / 2 // parent
		if i == j || !(heap[j].CmpKey() > heap[i].CmpKey()) {
			break
		}
		heap[i], heap[j] = heap[j], heap[i]
		j = i
	}
}
