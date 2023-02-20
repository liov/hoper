package leetcode

/**
378. 有序矩阵中第 K 小的元素

给定一个 n x n 矩阵，其中每行和每列元素均按升序排序，找到矩阵中第 k 小的元素。
请注意，它是排序后的第 k 小元素，而不是第 k 个不同的元素。



示例：

matrix = [
[ 1,  5,  9],
[10, 11, 13],
[12, 13, 15]
],
k = 8,

返回 13。


提示：
你可以假设 k 的值永远是有效的，1 ≤ k ≤ n^2 。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/kth-smallest-element-in-a-sorted-matrix
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun kthSmallest(matrix: Array<IntArray>, k: Int): Int {
  val n = matrix.size
  var left = matrix[0][0]
  var right = matrix[n - 1][n - 1]
  while (left < right) {
    val mid = right + left shr 1
    if (check(matrix, mid, k, n)) right = mid else  left = mid + 1
  }
  return left
}

fun check(matrix: Array<IntArray>, mid: Int, k: Int, n: Int): Boolean {
  var i = n - 1
  var j = 0
  var num = 0
  while (i >= 0 && j < n) {
    if (matrix[i][j] <= mid) {
      //每一列都比中值小
      num += i + 1
      j++
    } else i--
  }
  //比中值小的个数跟K比，大于K在左边，小于在右边
  return num >= k
}

