package middle

/*
78. 子集

中等

给你一个整数数组 nums ，数组中的元素 互不相同 。返回该数组所有可能的
子集
（幂集）。

解集 不能 包含重复的子集。你可以按 任意顺序 返回解集。

示例 1：

输入：nums = [1,2,3]
输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
示例 2：

输入：nums = [0]
输出：[[],[0]]

提示：

1 <= nums.length <= 10
-10 <= nums[i] <= 10
nums 中的所有元素 互不相同
*/
func subsets(nums []int) [][]int {
	n := len(nums)
	var ans [][]int
	tmp := make([]int, 0, len(nums))
	var backtrack func(k int)
	backtrack = func(k int) {
		if k == n {
			ans = append(ans, append([]int{}, tmp...))
			return
		}
		tmp = append(tmp, nums[k])
		backtrack(k + 1)
		tmp = tmp[:len(tmp)-1]
		backtrack(k + 1)

	}
	backtrack(0)
	return ans
}
