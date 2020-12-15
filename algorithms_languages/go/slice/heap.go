package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"
)

type MaxHeap []int

func (h *MaxHeap) Max() int             { return (*h)[0] }
func (h *MaxHeap) Len() int             { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool   { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)        { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *MaxHeap) Push(x interface{})   { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() (i interface{}) { i, *h = (*h)[len(*h)-1], (*h)[:len(*h)-1]; return i }

func main() {
	rand.Seed(time.Now().UnixNano())
	var h = new(MaxHeap)
	for i := 0; i < 10; i++ {
		v := rand.Intn(1000000)
		heap.Push(h, v)
	}
	fmt.Println(h)
	for i := 0; i < 10000; i++ {
		v := rand.Intn(1000000)
		if v < h.Max() {
			(*h)[0] = v
			heap.Fix(h, 0)
		}

	}
	fmt.Println(h)
	heap.Push(h, 10000) //推入，向上浮动
	fmt.Println(h)
	var h2 = MinHeap(make([]int, 10))
	for i := 0; i < 10000; i++ {
		v := rand.Intn(1000000)
		if v > h2.Min() {
			h2[0] = v
			heap.Fix(&h2, 0)
		}

	}
	fmt.Println(h2)
}

type MinHeap []int

func (h *MinHeap) Min() int             { return (*h)[0] }
func (h *MinHeap) Len() int             { return len(*h) }
func (h *MinHeap) Less(i, j int) bool   { return (*h)[i] < (*h)[j] }
func (h *MinHeap) Swap(i, j int)        { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *MinHeap) Push(x interface{})   { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() (i interface{}) { i, *h = (*h)[len(*h)-1], (*h)[:len(*h)-1]; return i }

func buildMaxHeap(a []int, heapSize int) {
	for i := heapSize / 2; i >= 0; i-- {
		maxHeapify(a, i, heapSize)
	}
}

func maxHeapify(a []int, i, heapSize int) {
	l, r, largest := i*2+1, i*2+2, i
	if l < heapSize && a[l] > a[largest] {
		largest = l
	}
	if r < heapSize && a[r] > a[largest] {
		largest = r
	}
	if largest != i {
		a[i], a[largest] = a[largest], a[i]
		maxHeapify(a, largest, heapSize)
	}
}
