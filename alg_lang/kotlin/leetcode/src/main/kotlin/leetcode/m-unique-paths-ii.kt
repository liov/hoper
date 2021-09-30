package leetcode

/**
不同路径 II

一个机器人位于一个 m x n 网格的左上角 （起始点在下图中标记为“Start” ）。

机器人每次只能向下或者向右移动一步。机器人试图达到网格的右下角（在下图中标记为“Finish”）。

现在考虑网格中有障碍物。那么从左上角到右下角将会有多少条不同的路径？


网格中的障碍物和空位置分别用 1 和 0 来表示。

说明：m 和 n 的值均不超过 100。

示例 1:

输入:
[
  [0,0,0],
  [0,1,0],
  [0,0,0]
]
输出: 2
解释:
3x3 网格的正中间有一个障碍物。
从左上角到右下角一共有 2 条不同的路径：
1. 向右 -> 向右 -> 向下 -> 向下
2. 向下 -> 向下 -> 向右 -> 向右

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/unique-paths-ii
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun uniquePathsWithObstacles(obstacleGrid: Array<IntArray>): Int {
  if (obstacleGrid.isEmpty() || obstacleGrid[0].isEmpty() || obstacleGrid[0][0] == 1 || obstacleGrid.last().last() == 1) return 0
  val m = obstacleGrid.size
  val n = obstacleGrid[0].size
  val mem = Array(m) { IntArray(n) }
  return uniquePathsDfs(obstacleGrid, 0, 0, mem)
}

fun uniquePathsDfs(obstacleGrid: Array<IntArray>, x: Int, y: Int, mem: Array<IntArray>): Int {
  if (x == obstacleGrid.size - 1 && y == obstacleGrid[0].size - 1) {
    mem[x][y] = 1
    return 1
  }
  if (x >= obstacleGrid.size || y >= obstacleGrid[0].size) return 0
  if (obstacleGrid[x][y] == 1) return 0
  if (mem[x][y] > 0) return mem[x][y] //遍历过直接返回
  val nextX = x + 1
  if (nextX < obstacleGrid.size) mem[x][y] += uniquePathsDfs(obstacleGrid, nextX, y, mem)
  val nextY = y + 1
  if (nextY < obstacleGrid[0].size) mem[x][y] += uniquePathsDfs(obstacleGrid, x, nextY, mem)
  return mem[x][y]
}

//动态规划
