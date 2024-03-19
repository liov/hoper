package leetcode

import java.lang.Math.pow

/**
343. 整数拆分

给定一个正整数 n，将其拆分为至少两个正整数的和，并使这些整数的乘积最大化。 返回你可以获得的最大乘积。

示例 1:

输入: 2
输出: 1
解释: 2 = 1 + 1, 1 × 1 = 1。
示例 2:

输入: 10
输出: 36
解释: 10 = 3 + 3 + 4, 3 × 3 × 4 = 36。
说明: 你可以假设 n 不小于 2 且不大于 58。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/integer-break
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun integerBreak(n: Int): Int {
  if (n==2) return 1
  if (n==3) return 2
  if (n==4) return 4
  var p = n/3
  var mod = n%3
  if (mod == 1) {
    mod = 4
    p--
  }
  if (mod == 0) mod = 1

  return myPow(3,p)*mod
}

fun myPow(x: Int, n: Int): Int {
  if (n == 0) return 1
  var ret = 1
  var vn = n
  var x1 = x
  while (vn > 0) {
    if (vn and 1 == 1) ret *= x1
    x1 *= x1
    vn = vn shr 1
  }
  return ret
}
