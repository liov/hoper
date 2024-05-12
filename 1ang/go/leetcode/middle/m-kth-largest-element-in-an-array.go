package middle

import "math"

/*

215. 数组中的第K个最大元素
已解答
中等

给定整数数组 nums 和整数 k，请返回数组中第 k 个最大的元素。

请注意，你需要找的是数组排序后的第 k 个最大的元素，而不是第 k 个不同的元素。

你必须设计并实现时间复杂度为 O(n) 的算法解决此问题。



示例 1:

输入: [3,2,1,5,6,4], k = 2
输出: 5
示例 2:

输入: [3,2,3,1,2,4,5,5,6], k = 4
输出: 4


提示：

1 <= k <= nums.length <= 105
-104 <= nums[i] <= 104
*/

func findKthLargest(nums []int, k int) int {
	heap := make([]int, k)
	for i := 0; i < k; i++ {
		heap[i] = math.MinInt
	}
	for _, num := range nums {
		if num > heap[0] {
			heap[0] = num
			down1(heap, 0)
		}
	}
	return heap[0]
}

func down1(nums []int, i int) {
	child := i*2 + 1
	for child < len(nums) {
		if child+1 < len(nums) && nums[child] > nums[child+1] {
			child++
		}
		if nums[child] > nums[i] {
			break
		}
		nums[i], nums[child] = nums[child], nums[i]
		i = child
		child = child*2 + 1
	}
}
