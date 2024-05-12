package leetcode

/**
201. 数字范围按位与

给定范围 [m, n]，其中 0 <= m <= n <= 2147483647，返回此范围内所有数字的按位与（包含 m, n 两端点）。

示例 1:

输入: [5,7]
输出: 4
示例 2:

输入: [0,1]
输出: 0

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/bitwise-and-of-numbers-range
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun rangeBitwiseAnd(m: Int, n: Int): Int {
  var n = n
  while (m < n)  n = n and (n - 1)
  return m and n
}

fun rangeBitwiseAndV2(m: Int, n: Int): Int {
  var m = m
  var n = n
  var shift = 0
  while(m<n){
    m = m shr 1
    n = n shr 1
    shift+=1
  }
  return m shl shift
}
