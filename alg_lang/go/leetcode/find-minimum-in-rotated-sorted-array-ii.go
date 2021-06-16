package leetcode

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func findMin(numbers []int) int {
	n := len(numbers)
	if n == 1 {
		return numbers[0]
	}
	if n == 2 {
		return min(numbers[0], numbers[1])
	}
	if numbers[0] < numbers[n-1] {
		return numbers[0]
	}
	if numbers[0] < numbers[n-1] {
		return numbers[0]
	}
	if numbers[n-1] < numbers[n-2] {
		return numbers[n-1]
	}
	var left = 0
	var right = n
	for left < right {
		mid := (left + right) >> 1
		if numbers[mid] > numbers[mid+1] {
			return numbers[mid+1]
		}
		if numbers[mid] < numbers[mid-1] {
			return numbers[mid]
		}
		if numbers[mid] < numbers[n-1] || numbers[mid] < numbers[0] {
			right = mid
		} else if numbers[mid] > numbers[n-1] || numbers[mid] > numbers[0] {
			left = mid
		} else {
			return min(findMin(numbers[0:mid]), findMin(numbers[mid:right]))
		}
	}
	return numbers[left]
}

func minArray(numbers []int) int {
	var low = 0
	var high = len(numbers) - 1
	for low < high {
		pivot := low + (high-low)/2
		switch {
		case numbers[pivot] < numbers[high]:
			high = pivot
		case numbers[pivot] > numbers[high]:
			low = pivot + 1
		default:
			high -= 1
		}
	}
	return numbers[low]
}
