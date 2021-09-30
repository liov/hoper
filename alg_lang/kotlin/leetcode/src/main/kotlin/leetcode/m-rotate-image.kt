package leetcode

/**
旋转图像

给定一个 n × n 的二维矩阵表示一个图像。

将图像顺时针旋转 90 度。

说明：

你必须在原地旋转图像，这意味着你需要直接修改输入的二维矩阵。请不要使用另一个矩阵来旋转图像。

示例 1:

给定 matrix =
[
[1,2,3],
[4,5,6],
[7,8,9]
],

原地旋转输入矩阵，使其变为:
[
[7,4,1],
[8,5,2],
[9,6,3]
]
示例 2:

给定 matrix =
[
[ 5, 1, 9,11],
[ 2, 4, 8,10],
[13, 3, 6, 7],
[15,14,12,16]
],

原地旋转输入矩阵，使其变为:
[
[15,13, 2, 5],
[14, 3, 4, 1],
[12, 6, 8, 9],
[16, 7,10,11]
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/rotate-image
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun rotate(matrix: Array<IntArray>): Unit {
  var n = matrix.size
  var x: Int
  var y: Int
  var tmp: Int
  while (n > 1) {
    x = (matrix.size - n) / 2 //当前阶的第一位
    y = x + n - 1 //当前阶的最后一位,matrix.size - 1 - x
    for (i in 0 until n - 1) {
      tmp = matrix[x][x + i]
      matrix[x][x + i] = matrix[y - i][x]
      matrix[y - i][x] = matrix[y][y - i]
      matrix[y][y - i] = matrix[x + i][y]
      matrix[x + i][y] = tmp
    }
    n -= 2
  }
}
