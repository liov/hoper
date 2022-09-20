package leetcode

/**
15. 三数之和

给你一个包含 n 个整数的数组nums，判断nums中是否存在三个元素 a，b，c ，使得a + b + c = 0 ？请你找出所有满足条件且不重复的三元组。

注意：答案中不可以包含重复的三元组。



示例：

给定数组 nums = [-1, 0, 1, 2, -1, -4]，

满足要求的三元组集合为：
[
[-1, 0, 1],
[-1, -1, 2]
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/3sum
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */

//516 ms , 在所有 Kotlin 提交中击败了 51.35% 的用户,这也太惨了
fun threeSum(nums: IntArray): List<List<Int>> {
  nums.sort()
  val ans = mutableListOf<List<Int>>() //val ans:List<List<Int>> 这种写法不能add，因为List没有add方法
  if (nums.isEmpty() || nums[0] > 0) return ans
  val map = HashMap<Int, Int>()
  for ((i, v) in nums.withIndex()) {
    if (v > 0) map[v] = i
  }

  var third: Int
  var zeroRecord = false //当0出现后，检查是否存在三个0的情况
  for (i in nums.indices) {
    if (i > 0 && nums[i] == nums[i - 1]) continue
    if (zeroRecord) break //0及0以上不用遍历
    if (nums[i] == 0) {
      if (nums.size - i > 2 && nums[i + 1] == 0 && nums[i + 2] == 0) ans.add(listOf(0, 0, 0))
      zeroRecord = true
    }
    for (j in i + 1 until nums.lastIndex) {
      third = 0 - nums[i] - nums[j]
      if (j > i + 1 && nums[j] == nums[j - 1]) continue //去重
      if (third <= 0) break //两数已经大于0不可能存在第三个数
      if (map[third] ?: 0 > j) ans.add(listOf(nums[i], nums[j], third)) //是否存在第三个数,第三个数只能在右边
    }
  }
  return ans
}

