package hard

/*
689. 三个无重叠子数组的最大和
给你一个整数数组 nums 和一个整数 k ，找出三个长度为 k 、互不重叠、且 3 * k 项的和最大的子数组，并返回这三个子数组。

以下标的数组形式返回结果，数组中的每一项分别指示每个子数组的起始位置（下标从 0 开始）。如果有多个结果，返回字典序最小的一个。



示例 1：

输入：nums = [1,2,1,2,6,7,5,1], k = 2
输出：[0,3,5]
解释：子数组 [1, 2], [2, 6], [7, 5] 对应的起始下标为 [0, 3, 5]。
也可以取 [2, 1], 但是结果 [1, 3, 5] 在字典序上更大。
示例 2：

输入：nums = [1,2,1,2,1,2,1,2,1], k = 2
输出：[0,2,4]


提示：

1 <= nums.length <= 2 * 10^4
1 <= nums[i] < 2^16
1 <= k <= floor(nums.length / 3)

https://leetcode-cn.com/problems/maximum-sum-of-3-non-overlapping-subarrays/
*/
// timeout...
func maxSumOfThreeSubarrays(nums []int, k int) []int {
	ans := make([]int, 3)
	arr := make([]int, len(nums)+1-k)
	for i := 0; i < k; i++ {
		arr[0] += nums[i]
	}
	for i := 1; i < len(arr); i++ {
		arr[i] = arr[i-1] + nums[i+k-1] - nums[i-1]
	}
	type Record struct {
		sum, x, y int
	}
	maxArr := make([]Record, len(arr)-k-k)
	var max1 int
	for j := k; j < len(nums)-k; j++ {
		for x := j + k; x < len(arr); x++ {
			tmp := arr[j] + arr[x]
			for y := 0; y <= j-k; y++ {
				if tmp > maxArr[y].sum {
					maxArr[y].sum = tmp
					maxArr[y].x, maxArr[y].y = j, x
				}
			}

		}
	}
	max1 = 0
	for i := 0; i < len(arr)-k*2; i++ {
		tmp := maxArr[i].sum + arr[i]
		if tmp > max1 {
			max1 = tmp
			ans[0], ans[1], ans[2] = i, maxArr[i].x, maxArr[i].y
		}
	}
	return ans
}

/*
我们使用三个大小为 k 的滑动窗口。设 sum1为第一个滑动窗口的元素和，该滑动窗口从 [0,k-1] 开始；sum2

为第二个滑动窗口的元素和，该滑动窗口从 [k,2k-1] 开始；sum3

	为第三个滑动窗口的元素和，该滑动窗口从 [2k,3k−1] 开始。

我们同时向右滑动这三个窗口，按照前言二的方法并维护 maxSum12

及其对应位置。每次滑动时，计算当前 maxSum12与 sum3

之和。统计这一过程中的 maxSum12+sum3的最大值及其对应位置。

对于题目要求的最小字典序，由于我们是从左向右遍历的，并且仅当元素和超过最大元素和时才修改最大元素和，从而保证求出来的下标列表是字典序最小的。

作者：LeetCode-Solution
链接：https://leetcode-cn.com/problems/maximum-sum-of-3-non-overlapping-subarrays/solution/san-ge-wu-zhong-die-zi-shu-zu-de-zui-da-4a8lb/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func maxSumOfThreeSubarrays2(nums []int, k int) (ans []int) {
	var sum1, maxSum1, maxSum1Idx int
	var sum2, maxSum12, maxSum12Idx1, maxSum12Idx2 int
	var sum3, maxTotal int
	for i := k * 2; i < len(nums); i++ {
		sum1 += nums[i-k*2]
		sum2 += nums[i-k]
		sum3 += nums[i]
		if i >= k*3-1 {
			if sum1 > maxSum1 {
				maxSum1 = sum1
				maxSum1Idx = i - k*3 + 1
			}
			if maxSum1+sum2 > maxSum12 {
				maxSum12 = maxSum1 + sum2
				maxSum12Idx1, maxSum12Idx2 = maxSum1Idx, i-k*2+1
			}
			if maxSum12+sum3 > maxTotal {
				maxTotal = maxSum12 + sum3
				ans = []int{maxSum12Idx1, maxSum12Idx2, i - k + 1}
			}
			sum1 -= nums[i-k*3+1]
			sum2 -= nums[i-k*2+1]
			sum3 -= nums[i-k+1]
		}
	}
	return
}

func maxSumOfThreeSubarrays3(nums []int, k int) []int {
	ans := []int{0, k, 2 * k}
	var sum1, sum2, sum3, sum12, sum int
	idx1, idx2, idx3 := 0, k, k
	for i := 0; i < k; i++ {
		sum1 += nums[i]
		sum2 += nums[i+k]
		sum3 += nums[i+2*k]
	}
	sum12 = sum1 + sum2
	sum = sum12 + sum3
	sum1tmp, sum2tmp, sum3tmp := sum1, sum2, sum3
	for i := 2*k + 1; i < len(nums)-k+1; i++ {
		sum1tmp += nums[i-k-1] - nums[i-2*k-1]
		if sum1tmp > sum1 {
			sum1 = sum1tmp
			idx1 = i - 2*k
		}
		sum2tmp += nums[i-1] - nums[i-k-1]
		if sum1+sum2tmp > sum12 {
			sum12 = sum1 + sum2tmp
			idx2, idx3 = idx1, i-k
		}

		sum3tmp += nums[i+k-1] - nums[i-1]
		if sum12+sum3tmp > sum {
			sum = sum12 + sum3tmp
			ans[0], ans[1], ans[2] = idx2, idx3, i
		}
	}
	return ans
}
