package leetcode

/**
 * 给定一个包含 n 个整数的数组 nums 和一个目标值 target，判断 nums 中是否存在四个元素 a，b，c 和 d ，使得 a + b + c + d 的值与 target 相等？找出所有满足条件且不重复的四元组。

注意：

答案中不可以包含重复的四元组。

示例：

给定数组 nums = [1, 0, -1, 0, -2, 2]，和 target = 0。

满足要求的四元组集合为：
[
[-1,  0, 0, 1],
[-2, -1, 1, 2],
[-2,  0, 0, 2]
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/4sum
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun fourSum(nums: IntArray, target: Int): List<List<Int>> {
  nums.sort()
  if (nums.isEmpty() || nums[0] > target) return mutableListOf()
  val map = HashMap<Int, Int>()
  for ((i, v) in nums.withIndex()) {
    if (v > 0) map[v] = i
  }
  return nSum(nums.toList(), target, 4, map)
}

fun nSum(nums: List<Int>, target: Int, n: Int, map: HashMap<Int, Int>): List<List<Int>> {
  val ans = mutableListOf<MutableList<Int>>()
  for (i in nums.indices) {
    if ((i > 0 && nums[i - 1] == nums[i]) || nums[i] > target) continue
    nSumIter(nums.subList(i + 1, nums.size), target - nums[i], n - 1, map)?.apply { add(nums[i]) }?.let { ans.add(it) }
  }
  return ans
}

fun nSumIter(nums: List<Int>, target: Int, n: Int, map: HashMap<Int, Int>): MutableList<Int>? {
  if (n != 2) {
    for (i in nums.indices) {
      if ((i > 0 && nums[i - 1] == nums[i]) || nums[i] > target) continue
      return nSumIter(nums.subList(i + 1, nums.size), target - nums[i], n - 1, map)?.apply { add(nums[i]) }
    }
  }

  for (i in nums.indices) {
    if (map[target - nums[i]] != null && map[target - nums[i]]!! > i + (map.size - nums.size)) {
      println("${target - nums[i]},${nums[i]}")
      return mutableListOf(target - nums[i], nums[i])
    }
    if ((i > 0 && nums[i - 1] == nums[i]) || nums[i] > target) continue
  }
  return null
}
