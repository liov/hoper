package leetcode

import java.util.*
import kotlin.collections.ArrayDeque
import kotlin.collections.ArrayList

/**
39. 组合总和

给定一个无重复元素的数组 candidates 和一个目标数 target ，找出 candidates 中所有可以使数字和为 target 的组合。

candidates 中的数字可以无限制重复被选取。

说明：

所有数字（包括 target）都是正整数。
解集不能包含重复的组合。
示例 1:

输入: candidates = [2,3,6,7], target = 7,
所求解集为:
[
[7],
[2,2,3]
]
示例 2:

输入: candidates = [2,3,5], target = 8,
所求解集为:
[
  [2,2,2,2],
  [2,3,3],
  [3,5]
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/combination-sum
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */

fun combinationSum(candidates: IntArray, target: Int): List<List<Int>> {
  candidates.sort()
  val set = HashSet<Int>()
  for (i in candidates.indices) {
    set.add(candidates[i])
  }
  val ans = ArrayList<List<Int>>()
  find(candidates, 0, target, ans, set, Stack())
  return ans
}

fun find(candidates: IntArray, start: Int, target: Int, ans: MutableList<List<Int>>, set: HashSet<Int>, path: Stack<Int>) {

  if (set.contains(target)) {
    path.push(target)
    ans.add(path.toList())
    path.pop()
    if (target == candidates[start]) return
  }

  for (i in start until candidates.size) {
    if (path.isNotEmpty() && candidates[i] < path.last()) continue
    path.push(candidates[i])
    if (target - candidates[i] < path.last()) {
      path.pop()
      break
    }
    var newStart = start
    if (path.size == 1) newStart = i
    find(candidates, newStart, target - candidates[i], ans, set, path)
    path.pop()
  }
}
