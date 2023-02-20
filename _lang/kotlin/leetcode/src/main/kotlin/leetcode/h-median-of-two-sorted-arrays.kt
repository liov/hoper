package leetcode

/**
4. 寻找两个正序数组的中位数

给定两个大小为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。

请你找出这两个正序数组的中位数，并且要求算法的时间复杂度为 O(log(m + n))。

你可以假设 nums1 和 nums2 不会同时为空。



示例 1:

nums1 = [1, 3]
nums2 = [2]

则中位数是 2.0
示例 2:

nums1 = [1, 2]
nums2 = [3, 4]

则中位数是 (2 + 3)/2 = 2.5

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/median-of-two-sorted-arrays
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
//效率不行 50%
fun findMedianSortedArrays(nums1: IntArray, nums2: IntArray): Double {
  val size = nums1.size + nums2.size;
  val even = size and 1 == 0 //偶数
  val halflen = size / 2
  var long = nums1 //不爽的就是不能一行声明多个
  var short = nums2
  if (nums2.size > nums1.size) {
    long = nums2
    short = nums1
  }

  var iMin = 0
  var iMax = short.size
  var i: Int //短数组的分隔位
  var j: Int //长数组的分隔位
  while (iMin <= iMax) {
    i = (iMin + iMax) / 2
    j = halflen - i
    when {
      i < iMax && short[i] < long[j - 1] -> iMin = i + 1
      i > iMin && short[i - 1] > long[j] -> iMax = i - 1
      else -> {
        val minRight = when {
          i == short.size -> long[j]
          j == long.size -> short[i]
          else -> kotlin.math.min(short[i], long[j])
        }
        if (!even) return minRight.toDouble()
        val maxLeft = when {
          i == 0 -> long[j - 1]
          j == 0 -> short[i - 1]
          else -> kotlin.math.max(short[i - 1], long[j - 1])
        }
        return (maxLeft + minRight) / 2.0
      }
    }
  }
  return 0.0
}

fun findMedianSortedArraysV2(nums1: IntArray, nums2: IntArray): Double {
  val size = nums1.size + nums2.size;
  val even = size and 1 == 0 //偶数
  val halflen = size / 2

  for ((i,j) in nums1.withIndex()){

  }
  return 0.0
}
