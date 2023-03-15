package slices

import "github.com/liov/hoper/server/go/lib/utils/def/constraints"

// TODO
func Sort[T comparable](slices []T) {

}

func quickSort[T constraints.Number](array []T, left, right int) {
	if left < right {
		x, i := array[right], left-1
		var temp T
		for j := left; j <= right; j++ {
			if array[j] <= x {
				i++
				temp = array[i]
				array[i] = array[j]
				array[j] = temp
			}
		}
		quickSort(array, left, i-1)
		quickSort(array, i+1, right)
	}
}

func bubbleSort[T constraints.Number](array []T) {
	low, high := 0, len(array)-1

	var tmp T
	var j int
	for low < high {
		for j = low; j < high; j++ {
			//正向冒泡,找到最大者
			if array[j] > array[j+1] {
				tmp = array[j]
				array[j] = array[j+1]
				array[j+1] = tmp
			}
			if array[low] > array[j] {
				tmp = array[low]
				array[low] = array[j]
				array[j] = tmp
			}
		}
		high--
		low++
	}
}
