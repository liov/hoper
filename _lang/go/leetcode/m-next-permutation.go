package leetcode

import "sort"

/*
31. 下一个排列
实现获取 下一个排列 的函数，算法需要将给定数字序列重新排列成字典序中下一个更大的排列（即，组合出下一个更大的整数）。

如果不存在下一个更大的排列，则将数字重新排列成最小的排列（即升序排列）。

必须 原地 修改，只允许使用额外常数空间。



示例 1：

输入：nums = [1,2,3]
输出：[1,3,2]
示例 2：

输入：nums = [3,2,1]
输出：[1,2,3]
示例 3：

输入：nums = [1,1,5]
输出：[1,5,1]
示例 4：

输入：nums = [1]
输出：[1]


提示：

1 <= nums.length <= 100
0 <= nums[i] <= 100

https://leetcode-cn.com/problems/next-permutation/
*/

func nextPermutation(nums []int) {

	for i := len(nums) - 1; i > 0; i-- {
		if nums[i] > nums[i-1] {
			minIdx := i
			for j := i; j < len(nums); j++ {
				if nums[j] <= nums[i-1] {
					nums[j-1], nums[i-1] = nums[i-1], nums[j-1]
					sort.Ints(nums[i:])
					return
				}
				if nums[j] < nums[minIdx] && nums[j] > nums[i-1] {
					minIdx = j
				}
			}
			nums[minIdx], nums[i-1] = nums[i-1], nums[minIdx]
			sort.Ints(nums[i:])
			return
		}
	}
	sort.Ints(nums)
}
