package leetcode

/**
给定一组不含重复元素的整数数组 nums，返回该数组所有可能的子集（幂集）。

说明：解集不能包含重复的子集。

示例:

输入: nums = [1,2,3]
输出:
[
[3],
  [1],
  [2],
  [1,2,3],
  [1,3],
  [2,3],
  [1,2],
  []
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/subsets
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun subsets(nums: IntArray): List<List<Int>> {
  if (nums.isEmpty()) return emptyList()
  val ret = ArrayList<ArrayList<Int>>()
  ret.add(ArrayList())
  for (i in nums.indices) {
    var start = 0
    while (start + i < nums.size) {
      val list = ArrayList<Int>()
      for (x in start until start + i) {
        list.add(nums[x])
      }
      for (x in start + i until nums.size) {
        val sublist = (list.clone() as ArrayList<Int>)
        sublist.add(nums[x])
        ret.add(sublist)
      }
      start++
      if (i == 0) break
    }
  }
  return ret
}
