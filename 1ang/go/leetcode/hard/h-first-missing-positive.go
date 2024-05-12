package hard

/*
41. 缺失的第一个正数

给你一个未排序的整数数组 nums ，请你找出其中没有出现的最小的正整数。

请你实现时间复杂度为 O(n) 并且只使用常数级别额外空间的解决方案。

示例 1：

输入：nums = [1,2,0]
输出：3
示例 2：

输入：nums = [3,4,-1,1]
输出：2
示例 3：

输入：nums = [7,8,9,11,12]
输出：1

提示：

1 <= nums.length <= 5 * 10^5
-2^31 <= nums[i] <= 2^31 - 1

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/first-missing-positive
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
func firstMissingPositive(nums []int) int {
	for i := range nums {

		if nums[i] > len(nums) || nums[i] <= 0 || nums[i] == i+1 {
			continue
		}
		idx := nums[i] - 1
	Lable:
		idx2 := nums[idx] - 1
		if nums[idx] == idx+1 {
			continue
		}
		nums[idx] = idx + 1
		if idx2 > len(nums)-1 || idx2 < 0 {
			continue
		}
		idx = idx2

		goto Lable
	}
	for i := range nums {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return len(nums) + 1
}
