package leetcode

/**
215. 数组中的第K个最大元素

在未排序的数组中找到第 k 个最大的元素。请注意，你需要找的是数组排序后的第 k 个最大的元素，而不是第 k 个不同的元素。

示例 1:

输入: [3,2,1,5,6,4] 和 k = 2
输出: 5
示例 2:

输入: [3,2,3,1,2,4,5,5,6] 和 k = 4
输出: 4
说明:

你可以假设 k 总是有效的，且 1 ≤ k ≤ 数组的长度。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/kth-largest-element-in-an-array
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun findKthLargestV1(nums: IntArray, k: Int): Int {
  var max = 0
  var num = 0
  while (num < k) {
    max = num
    for (i in num until nums.size) {
      if (nums[i] > nums[max]) max = i
    }
    nums.swap(max, num)
    num++
  }
  return nums[k - 1]
}

fun findKthLargest(nums: IntArray, k: Int): Int {
  findTopMaxNInPlace(nums, k)
  return nums[0]
}
