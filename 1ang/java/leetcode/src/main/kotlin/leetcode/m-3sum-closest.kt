package leetcode

/**

16. 最接近的三数之和
给定一个包括 n 个整数的数组 nums 和 一个目标值 target。找出 nums 中的三个整数，使得它们的和与 target 最接近。返回这三个数的和。假定每组输入只存在唯一答案。



示例：

输入：nums = [-1,2,1,-4], target = 1
输出：2
解释：与 target 最接近的和是 2 (-1 + 2 + 1 = 2) 。


提示：

3 <= nums.length <= 10^3
-10^3 <= nums[i] <= 10^3
-10^4 <= target <= 10^4
 */
fun threeSumClosest(nums: IntArray, target: Int): Int {
  if (nums.size == 3) return nums[0] + nums[1] + nums[2]
  nums.sort()
  if (nums[0] * 3 >= target || target <= -3000) return nums[0] + nums[1] + nums[2]
  if (nums.last() * 3 <= target || target >= 3000) return nums[nums.size - 1] + nums[nums.size - 2] + nums[nums.size - 3]
  var ret = 10000
  var sum = nums[0] + nums[1] + nums[2]
  for (i in nums.indices) {
    var left = i + 1
    var right = nums.size - 1
    while (left < right) {
      sum = nums[i] + nums[left] + nums[right]
      //可以优化的
      if (kotlin.math.abs(target - ret) > kotlin.math.abs(target - sum)) ret = sum
      when {
        sum > target -> right--
        sum < target -> left++
        sum == target -> return target
      }
    }
  }
  return ret
}
