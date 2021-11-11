package leetcode

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for times := rand.Intn(10); times >= 0; times-- {
		size := rand.Intn(20)
		arr, arr1 := make([]int, size), make([]int, size)
		for i := range arr {
			num := rand.Intn(100000)
			arr[i], arr1[i] = num, num
		}
		res, res1 := NewMaxHeapFromArr(arr), NewMaxHeapFromArr2(arr1)
		fmt.Println(res, res1)
		fmt.Println(res.check(), res1.check())
	}
}
