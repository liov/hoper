package leetcode

import (
	"math"
)

type MaxHeap []int

func NewMaxHeap(l int) MaxHeap {
	maxHeap := make(MaxHeap, l)
	for i := range maxHeap {
		maxHeap[i] = math.MaxInt
	}
	return maxHeap
}

func NewMaxHeapFromArr(arr []int) MaxHeap {
	heap := MaxHeap(arr)
	for i := 1; i < len(arr); i++ {
		heap.adjustUp(i)
	}
	return heap
}

func NewMaxHeapFromArr2(arr []int) MaxHeap {
	heap := MaxHeap(arr)
	for i := len(arr)/2 - 1; i >= 0; i-- {
		heap.adjustDown(i)
	}
	return heap
}

func (heap MaxHeap) Put(val int) {
	if val > heap[0] {
		return
	}
	heap[0] = val
	heap.adjustDown(0)
}

func (heap MaxHeap) adjustUp(i int) {
	p := parent(i)
	for p >= 0 && heap[i] > heap[p] {
		Swap(heap, i, p)
		i = p
		p = parent(i)
	}

}

func (heap MaxHeap) adjustDown(i int) {
	child := leftChild(i)
	for child < len(heap) {
		if child+1 < len(heap) && heap[child+1] > heap[child] {
			child++
		}
		if heap[i] >= heap[child] {
			break
		}
		Swap(heap, i, child)
		i = child
		child = leftChild(i)
	}
}

func (heap MaxHeap) check() bool {
	for i := range heap {
		l, r := i*2+1, i*2+2
		if l < len(heap) && heap[i] < heap[l] {
			return false
		}
		if r < len(heap) && heap[i] < heap[r] {
			return false
		}
	}
	return true
}
