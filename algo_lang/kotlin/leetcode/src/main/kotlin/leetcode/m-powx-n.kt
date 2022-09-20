package leetcode

/**
50. Pow(x, n)

实现 pow(x, n) ，即计算 x 的 n 次幂函数。

示例 1:

输入: 2.00000, 10
输出: 1024.00000
示例 2:

输入: 2.10000, 3
输出: 9.26100
示例 3:

输入: 2.00000, -2
输出: 0.25000
解释: 2^-2 = 1/2^2 = 1/4 = 0.25
说明:

-100.0 < x < 100.0
n 是 32 位有符号整数，其数值范围是 [−2^31, 2^31 − 1] 。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/powx-n
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */

fun myPow(x: Double, n: Int): Double {
  if (n == 0) return 1.0
  var ret = 1.0
  var vn = if (n < 0) -n else n
  if (n == Int.MIN_VALUE) vn = Int.MAX_VALUE
  var x1 = x
  while (vn > 0) {
    //n为15x^(1111),x^(1110)*x,x^(2*(111))*x,x^(2*(110))*x^2
    //最后一位为1，则拆解开，剩余部分快速幂
    if (vn and 1 == 1) ret *= x1
    x1 *= x1
    vn = vn shr 1
  }
  if (n == Int.MIN_VALUE) ret*=x
  return if (n > 0) ret else 1 / ret
}
