package maxheap

import (
	"golang.org/x/exp/constraints"
)

// 最大（小）堆是指在树中，存在一个结点而且该结点有儿子结点，该结点的data域值都不小于（大于）其儿子结点的data域值

// 最大堆 可用于保留前n个最小元素
type MaxHeap[T constraints.Ordered] []T

func New[T constraints.Ordered](l int) MaxHeap[T] {
	maxHeap := make(MaxHeap[T], 0, l)
	return maxHeap
}

func NewFromArray[T constraints.Ordered](arr []T) MaxHeap[T] {
	heap := MaxHeap[T](arr)
	for i := 1; i < len(arr); i++ {
		heap.adjustUp(i)
	}
	return arr
}

func (heap MaxHeap[T]) Put(val T) {
	if len(heap) < cap(heap) {
		heap = append(heap, val)
		for i := 1; i < len(heap); i++ {
			heap.adjustUp(i)
		}
		return
	}
	if val > heap[0] {
		return
	}
	heap[0] = val
	heap.adjustDown(0)
}

func (heap MaxHeap[T]) swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}

func parent(i int) int {
	return (i - 1) / 2
}
func leftChild(i int) int {
	return i*2 + 1
}

func (heap MaxHeap[T]) adjustUp(i int) {
	p := parent(i)
	for p >= 0 && heap[i] > heap[p] {
		heap.swap(i, p)
		i = p
		p = parent(i)
	}

}

func (heap MaxHeap[T]) adjustDown(i int) {
	child := leftChild(i)
	for child < len(heap) {
		if child+1 < len(heap) && heap[child+1] > heap[child] {
			child++
		}
		if heap[i] >= heap[child] {
			break
		}
		heap.swap(i, child)
		i = child
		child = leftChild(i)
	}
}
