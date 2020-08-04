package leetcode

/**
 * 你这个学期必须选修 numCourse 门课程，记为 0 到 numCourse-1 。

在选修某些课程之前需要一些先修课程。 例如，想要学习课程 0 ，你需要先完成课程 1 ，我们用一个匹配来表示他们：[0,1]

给定课程总量以及它们的先决条件，请你判断是否可能完成所有课程的学习？

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/course-schedule
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
//一遍过，哈哈哈
//执行用时： 1336 ms , 在所有 Kotlin 提交中击败了 7.69% 的用户
//内存消耗： 54.9 MB , 在所有 Kotlin 提交中击败了 100.00% 的用户
//java 版本467ms，难道kotlin真比java慢3倍
fun canFinish(numCourses: Int, prerequisites: Array<IntArray>): Boolean {
  val map = HashMap<Int,HashSet<Int>>()
  for (arr in prerequisites){
    if(map.containsKey(arr[1])) map[arr[1]]!!.add(arr[0])
    else map[arr[1]] = hashSetOf(arr[0])
  }
  val set = HashSet<Int>()
  val mem = HashSet<Int>()
  for ((k,_) in map){
    if (!dfs(k,set,mem,map)) return false
  }
  return true
}

fun dfs(num:Int,set:HashSet<Int>,mem:HashSet<Int>,map: HashMap<Int,HashSet<Int>>):Boolean{
  if(set.contains(num)) return false
  if(mem.contains(num)) return true
  set.add(num)
  if(map.containsKey(num)){
    for(n in map[num]!!){
      if(mem.contains(n)) return true
      if (!dfs(n,set,mem,map)) return false
    }
  }
  set.remove(num)
  return true
}
