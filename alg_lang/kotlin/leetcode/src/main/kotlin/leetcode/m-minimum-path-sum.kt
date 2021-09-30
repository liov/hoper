package leetcode

/**
64. 最小路径和

给定一个包含非负整数的 m x n 网格，请找出一条从左上角到右下角的路径，使得路径上的数字总和为最小。

说明：每次只能向下或者向右移动一步。

示例:

输入:
[
[1,3,1],
[1,5,1],
[4,2,1]
]
输出: 7
解释: 因为路径 1→3→1→1→1 的总和最小。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/minimum-path-sum
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun minPathSum(grid: Array<IntArray>): Int {
  if (grid.isEmpty()) return 0
  val dp = IntArray(grid[0].size)

  for (x in grid.size - 1 downTo 0) {
    for (y in grid[0].size - 1 downTo 0) {
      if (x == grid.size - 1) {
        //原本应该首先单独判断的，放进来了，代码会简洁，可读性会变差
        if (y == grid[0].size - 1) dp[y] = grid[grid.size - 1][grid[0].size - 1]
        else dp[y] = dp[y + 1] + grid[x][y]
        continue
      }
      //原本应该分开写的
      if (y == grid[0].size - 1) dp[y] = dp[y] + grid[x][y]
      else dp[y] = kotlin.math.min(dp[y], dp[y + 1]) + grid[x][y]
    }
  }
  return dp[0]
}

data class Point(val x: Int, val y: Int)

//dfs超时
fun minPathSumV2(grid: Array<IntArray>): Int {
  if (grid.isEmpty()) return 0
  val map = HashMap<Point, Int>()
  return dfs(Point(0, 0), map, grid, 0)
}

fun dfs(p: Point, map: HashMap<Point, Int>, grid: Array<IntArray>, sum: Int): Int {
  if (map[p] != null) return map[p]!!
  if (p.x == grid.size - 1) {
    val ret = if (p.y == grid[0].size - 1) sum + grid[p.x][p.y]
    else dfs(Point(p.x, p.y + 1), map, grid, sum) + grid[p.x][p.y]
    map[p] = ret
    return ret
  }
  val ret =  if (p.y == grid[0].size - 1) dfs(Point(p.x + 1, p.y), map, grid, sum) + grid[p.x][p.y]
  else kotlin.math.min(dfs(Point(p.x + 1, p.y), map, grid, sum), dfs(Point(p.x, p.y + 1), map, grid, sum)) + grid[p.x][p.y]
  map[p] = ret
  return ret
}
