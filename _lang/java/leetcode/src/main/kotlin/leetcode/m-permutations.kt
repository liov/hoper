package leetcode

/**
46. 全排列

给定一个 没有重复 数字的序列，返回其所有可能的全排列。

示例:

输入: [1,2,3]
输出:
[
[1,2,3],
[1,3,2],
[2,1,3],
[2,3,1],
[3,1,2],
[3,2,1]
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/permutations
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */


fun permute(nums: IntArray): List<List<Int>> {
  if (nums.size == 1) return listOf(listOf(nums[0]))
  var count = 1
  for (i in 2..nums.size) {
    count *= i
  }
  val ret = ArrayList<ArrayList<Int>>(count)
  ret.add(ArrayList<Int>(nums.size).apply { this.add(nums[0]) })
  for (i in 1 until nums.size) {
    for (j in 0 until ret.size) {
      for (x in 0 until ret[j].size) {
        ret.add((ret[j].clone() as ArrayList<Int>).apply { this.add(x, nums[i]) } )
      }
      ret[j].add(nums[i])
    }
  }
  return ret
}
