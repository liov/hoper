package hard

/*
4. 寻找两个正序数组的中位数
困难

给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 中位数 。

算法的时间复杂度应该为 O(log (m+n)) 。



示例 1：

输入：nums1 = [1,3], nums2 = [2]
输出：2.00000
解释：合并数组 = [1,2,3] ，中位数 2
示例 2：

输入：nums1 = [1,2], nums2 = [3,4]
输出：2.50000
解释：合并数组 = [1,2,3,4] ，中位数 (2 + 3) / 2 = 2.5




提示：

nums1.length == m
nums2.length == n
0 <= m <= 1000
0 <= n <= 1000
1 <= m + n <= 2000
-106 <= nums1[i], nums2[i] <= 106
*/

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {

	short, long := nums1, nums2
	if len(nums1) > len(nums2) {
		short, long = nums2, nums1
	}

	shortLen, longLen := len(short), len(long)
	even := (shortLen+longLen)&1 == 0
	middle := (shortLen + longLen) >> 1
	var i, j int
	iMin, iMax := 0, shortLen
	for iMin <= iMax {
		i = (iMin + iMax) >> 1
		j = middle - i
		if i < iMax && short[i] < long[j-1] {
			iMin++
		} else if i > iMin && short[i-1] > long[j] {
			iMax--
		} else {
			var maxRight, maxLeft int
			if i == shortLen {
				maxRight = long[j]
			} else if j == longLen {
				maxRight = short[i]
			} else {
				maxRight = min(short[i], long[j])
			}
			if !even {
				return float64(maxRight)
			}
			if i == 0 {
				maxLeft = long[j-1]
			} else if j == 0 {
				maxLeft = short[i-1]
			} else {
				maxLeft = max(short[i-1], long[j-1])
			}
			return float64(maxLeft+maxRight) / 2.0
		}
	}
	return 0.0
}
