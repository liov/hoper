package middle

import "math"

/*
45. 跳跃游戏 II
中等

给定一个长度为 n 的 0 索引整数数组 nums。初始位置为 nums[0]。

每个元素 nums[i] 表示从索引 i 向前跳转的最大长度。换句话说，如果你在 nums[i] 处，你可以跳转到任意 nums[i + j] 处:

0 <= j <= nums[i]
i + j < n
返回到达 nums[n - 1] 的最小跳跃次数。生成的测试用例可以到达 nums[n - 1]。

示例 1:

输入: nums = [2,3,1,1,4]
输出: 2
解释: 跳到最后一个位置的最小跳跃数是 2。

	从下标为 0 跳到下标为 1 的位置，跳 1 步，然后跳 3 步到达数组的最后一个位置。

示例 2:

输入: nums = [2,3,0,1,4]
输出: 2

提示:

1 <= nums.length <= 104
0 <= nums[i] <= 1000
题目保证可以到达 nums[n-1]
*/
// 超时版本
func jump(nums []int) int {
	j := nums[0]
	if j == 0 || len(nums) == 1 {
		return 0
	}
	if j >= len(nums)-1 {
		return 1
	}
	minTimes := math.MaxInt
	for ; j > 0; j-- {
		times := jump(nums[j:])
		if times > 0 {
			minTimes = min(times, minTimes)
		}
	}
	return minTimes + 1
}

func jump2(nums []int) int {
	if len(nums) == 1 {
		return 0
	}
	most, latMost, steps := 0, 0, 0
	for i, num := range nums {

		if i <= most {
			most = max(most, num+i)
		}
		if most >= len(nums)-1 {
			return steps + 1
		}
		if i == latMost {
			latMost = most
			steps++
		}

	}
	return steps
}

func jump3(nums []int) int {
	length := len(nums)
	maxPosition, end, steps := 0, 0, 0
	for i := 0; i < length-1; i++ {
		maxPosition = max(maxPosition, i+nums[i])
		if i == end {
			end = maxPosition
			steps++
		}
	}
	return steps
}
