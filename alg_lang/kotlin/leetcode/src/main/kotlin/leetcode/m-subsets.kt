package leetcode

import java.util.*


/**
子集

给定一组不含重复元素的整数数组 nums，返回该数组所有可能的子集（幂集）。

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
fun subsetsV2(nums: IntArray): List<List<Int>> {
  if (nums.isEmpty()) return emptyList()
  val ret = ArrayList<ArrayList<Int>>()
  for (i in 0..nums.size) {
    backtrack(0,i,ArrayList(),ret,nums)
  }
  return ret
}

fun backtrack(first: Int,n:Int,curr:ArrayList<Int>,ret:ArrayList<ArrayList<Int>>,nums: IntArray){
  if(curr.size == n) {
    ret.add(ArrayList(curr))
    return
  }
  for (i in first until nums.size){
    curr.add(nums[i])
    backtrack(i+1,n,curr,ret,nums)
    //kt大坑
    curr.removeAt(curr.lastIndex)
  }
}

fun subsets(nums: IntArray): List<List<Int>> {
  val res: MutableList<List<Int>> = ArrayList()
  res.add(ArrayList())
  for (i in nums.indices) {
    val all = res.size
    for (j in 0 until all) {
      val tmp = ArrayList(res[j])
      tmp.add(nums[i])
      res.add(tmp)
    }
  }
  return res
}
