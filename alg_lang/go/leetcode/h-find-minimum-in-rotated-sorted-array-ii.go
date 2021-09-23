package leetcode

/**
154. 寻找旋转排序数组中的最小值 II
已知一个长度为 n 的数组，预先按照升序排列，经由 1 到 n 次 旋转 后，得到输入数组。例如，原数组 nums = [0,1,4,4,5,6,7] 在变化后可能得到：
若旋转 4 次，则可以得到 [4,5,6,7,0,1,4]
若旋转 7 次，则可以得到 [0,1,4,4,5,6,7]
注意，数组 [a[0], a[1], a[2], ..., a[n-1]] 旋转一次 的结果为数组 [a[n-1], a[0], a[1], a[2], ..., a[n-2]] 。

给你一个可能存在 重复 元素值的数组 nums ，它原来是一个升序排列的数组，并按上述情形进行了多次旋转。请你找出并返回数组中的 最小元素 。



示例 1：

输入：nums = [1,3,5]
输出：1
示例 2：

输入：nums = [2,2,2,0,1]
输出：0


提示：

n == nums.length
1 <= n <= 5000
-5000 <= nums[i] <= 5000
nums 原来是一个升序排序的数组，并进行了 1 至 n 次旋转


进阶：

这道题是 寻找旋转排序数组中的最小值 的延伸题目。
允许重复会影响算法的时间复杂度吗？会如何影响，为什么？
*/

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
