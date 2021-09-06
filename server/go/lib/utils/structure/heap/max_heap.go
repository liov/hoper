package heap

import "math"

type MaxHeap []int

func New(l int) MaxHeap {
	maxHeap := make(MaxHeap, l)
	for i := range maxHeap {
		maxHeap[i] = math.MaxInt
	}
	return maxHeap
}

func NewFromArr(arr []int) MaxHeap {
	for i := 1; i < len(arr); i++ {
		adjustUp(arr, i)
	}
	return arr
}

func (heap MaxHeap) Put(val int) {
	adjustDown(heap, val)
}

func swap(heap []int, i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}

func parent(i int) int {
	return (i - 1) / 2
}
func leftChild(i int) int {
	return i*2 + 1
}

func adjustUp(heap []int, i int) {
	p := parent(i)
	for p >= 0 && heap[i] > heap[p] {
		swap(heap, i, p)
		i = p
		p = parent(i)
	}

}

func adjustDown(heap []int, v int) {
	if v > heap[0] {
		return
	}
	heap[0] = v
	i := 0
	child := leftChild(0)
	for child < len(heap) {
		if child+1 < len(heap) && heap[child+1] > heap[child] {
			child++
		}
		if heap[i] >= heap[child] {
			break
		}
		swap(heap, i, child)
		i = child
		child = leftChild(i)
	}
}
