package leetcode

/**
415. 字符串相加

给定两个字符串形式的非负整数 num1 和num2 ，计算它们的和。

注意：

num1 和num2 的长度都小于 5100.
num1 和num2 都只包含数字 0-9.
num1 和num2 都不包含任何前导零。
你不能使用任何內建 BigInteger 库， 也不能直接将输入的字符串转换为整数形式。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/add-strings
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun addStrings(num1: String, num2: String): String {
  if(num1 == "0") return num2
  if(num2 == "0") return num1
  var short = num1
  var long = num2
  if (num1.length > num2.length) {
    short = num2
    long = num1
  }
  val m = short.length - 1
  val n = long.length - 1
  var carry = 0
  val ret = StringBuilder()
  for (i in 0..n) {
    val sum = (if (i <= m) (short[m - i] - '0') else 0) + (long[n - i] - '0') + carry
    if (sum >= 10) {
      carry = 1
      ret.append('0' + sum - 10)
    } else {
      carry = 0
      ret.append('0' + sum)
    }
  }
  if(carry == 1) ret.append('1')
  return ret.reverse().toString()
}
