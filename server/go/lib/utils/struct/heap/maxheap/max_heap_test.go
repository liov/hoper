package maxheap

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMaxHeap(t *testing.T) {
	var arr []int
	for i := 0; i < 10; i++ {
		arr = append(arr, rand.Intn(10000))
	}
	maxHeap := NewFromArr(arr)
	fmt.Println(maxHeap)
	for i := 0; i < 10; i++ {
		maxHeap.Put(rand.Intn(10000))
	}
	fmt.Println(maxHeap)
}
