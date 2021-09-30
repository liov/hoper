package leetcode

//只出现一次的数字 https://leetcode-cn.com/problems/single-number/
fun singleNumber(nums: IntArray): Int {
  return nums.reduce { x, y -> x xor y }
}
