package middle

/*
229. 求众数 II
给定一个大小为 n 的整数数组，找出其中所有出现超过 ⌊ n/3 ⌋ 次的元素。

示例 1：

输入：[3,2,3]
输出：[3]
示例 2：

输入：nums = [1]
输出：[1]
示例 3：

输入：[1,1,1,3,3,2,2,2]
输出：[1,2]

提示：

1 <= nums.length <= 5 * 10^4
-10^9 <= nums[i] <= 10^9

进阶：尝试设计时间复杂度为 O(n)、空间复杂度为 O(1)的算法解决此问题。

https://leetcode-cn.com/problems/majority-element-ii/
*/
func majorityElement(nums []int) []int {
	var ret []int
	times := len(nums) / 3
	m := make(map[int]int)
	for i := range nums {
		num := m[nums[i]]
		num++
		m[nums[i]] = num
		if num > times {
			ret = append(ret, nums[i])
			m[nums[i]] = -len(nums)
		}

	}
	return ret
}

func majorityElement2(nums []int) []int {
	var ans []int
	element1, element2 := 0, 0
	vote1, vote2 := 0, 0

	for _, num := range nums {
		if vote1 > 0 && num == element1 { // 如果该元素为第一个元素，则计数加1
			vote1++
		} else if vote2 > 0 && num == element2 { // 如果该元素为第二个元素，则计数加1
			vote2++
		} else if vote1 == 0 { // 选择第一个元素
			element1 = num
			vote1++
		} else if vote2 == 0 { // 选择第二个元素
			element2 = num
			vote2++
		} else { // 如果三个元素均不相同，则相互抵消1次
			vote1--
			vote2--
		}
	}

	cnt1, cnt2 := 0, 0
	for _, num := range nums {
		if vote1 > 0 && num == element1 {
			cnt1++
		}
		if vote2 > 0 && num == element2 {
			cnt2++
		}
	}

	// 检测元素出现的次数是否满足要求
	if vote1 > 0 && cnt1 > len(nums)/3 {
		ans = append(ans, element1)
	}
	if vote2 > 0 && cnt2 > len(nums)/3 {
		ans = append(ans, element2)
	}
	return ans
}
