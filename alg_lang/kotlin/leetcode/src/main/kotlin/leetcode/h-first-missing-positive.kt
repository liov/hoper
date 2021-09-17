package leetcode

/**
41. 缺失的第一个正数

给你一个未排序的整数数组，请你找出其中没有出现的最小的正整数。



示例 1:

输入: [1,2,0]
输出: 3
示例 2:

输入: [3,4,-1,1]
输出: 2
示例 3:

输入: [7,8,9,11,12]
输出: 1


提示：

你的算法的时间复杂度应为O(n)，并且只能使用常数级别的额外空间。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/first-missing-positive
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun firstMissingPositive(nums: IntArray): Int {
  if (nums.isEmpty()) return 1
  var minZeroIdx = 0
  var j: Int
  var tmp: Int
  for (i in nums.indices) {
    j = nums[i]
    if (j == i + 1) continue
    while (j in 1..nums.size && nums[j - 1] != j) {
      tmp = nums[j - 1]
      nums[j - 1] = j
      j = tmp
    }
    if (j !in 1..nums.size || nums[i] != i+1) {
      nums[i] = 0
      if (nums[minZeroIdx] != 0) minZeroIdx = i
    }
  }
  for(i in nums.indices){
    if (nums[i] == 0) return i+1
  }
  return if (nums[0] == 1 && minZeroIdx == 0 && nums.last() == nums.size) return nums.size + 1 else 1
}
