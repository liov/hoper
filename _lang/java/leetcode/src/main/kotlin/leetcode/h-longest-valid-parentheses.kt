package leetcode

import java.util.*
import kotlin.collections.ArrayList


/**
32. 最长有效括号

给定一个只包含 '(' 和 ')' 的字符串，找出最长的包含有效括号的子串的长度。

示例 1:

输入: "(()"
输出: 2
解释: 最长有效括号子串为 "()"
示例 2:

输入: ")()())"
输出: 4
解释: 最长有效括号子串为 "()()"

https://leetcode-cn.com/problems/longest-valid-parentheses/
 */
fun longestValidParenthesesV1(s: String): Int {
  if (s.length < 2) return 0
  val stack = Stack<CharAndIndex>()
  val list = ArrayList<Int>()
  var len = 0
  for (i in s.indices) {
    if (stack.size > 0 && stack.last().value == '(' && s[i] == ')') {
      stack.pop()
      list.add(i)
    } else stack.push(CharAndIndex(s[i], i))
  }
  //为空说明整串符合
  if (stack.isEmpty()) return s.length
  if (list.isEmpty()) return 0
  stack.push(CharAndIndex('(', s.length))
  var max = 0
  var j = 0
  loop@ for (i in stack.indices) {
    while (list[j] < stack[i].index) {
      len += 2
      j++
      if (j == list.size) {
        max = kotlin.math.max(max, len)
        break@loop
      }
    }
    max = kotlin.math.max(max, len)
    len = 0
  }
  return max
}

class CharAndIndex(val value: Char, val index: Int)

@JvmInline
value class CharIndex(val value: Int)

/**
在这种方法中，我们利用两个计数器 left 和 right 。
首先，我们从左到右遍历字符串，对于遇到的每个 (，我们增加 left 计数器，对于遇到的每个 ) ，我们增加 right 计数器。
每当 left 计数器与 right 计数器相等时，我们计算当前有效字符串的长度，并且记录目前为止找到的最长子字符串。
如果 right 计数器比 left 计数器大时，我们将 left 和 right 计数器同时变回 0 。

接下来，我们从右到左做一遍类似的工作。

 */
//比我的还慢
fun longestValidParenthesesV2(s: String): Int {
  var left = 0
  var right = 0
  var maxlength = 0
  for (element in s) {
    if (element == '(') left++ else right++
    if (left == right) maxlength = kotlin.math.max(maxlength, 2 * right)
    else if (right >= left) {
      right = 0
      left = 0
    }
  }
  right = 0
  left = 0
  for (i in s.length - 1 downTo 0) {
    if (s[i] == '(') left++ else right++
    if (left == right) maxlength = kotlin.math.max(maxlength, 2 * left)
    else if (left >= right) {
      right = 0
      left = right
    }
  }
  return maxlength
}

/**
这个问题可以通过动态规划解决。我们定义一个 dp 数组，其中第 i 个元素表示以下标为 i 的字符结尾的最长有效子字符串的长度。我们将 dp 数组全部初始化为 0 。现在，很明显有效的子字符串一定以 ')' 结尾。
这进一步可以得出结论：以 '(' 结尾的子字符串对应的 dp 数组位置上的值必定为 0 。所以说我们只需要更新 ')' 在 dp 数组中对应位置的值。

为了求出 dp 数组，我们每两个字符检查一次，如果满足如下条件

s[i] =')' 且 s[i - 1] = '(' ，也就是字符串形如 "……()"，我们可以推出：

dp[i]=dp[i−2]+2

我们可以进行这样的转移，是因为结束部分的 "()" 是一个有效子字符串，并且将之前有效子字符串的长度增加了 2 。

s[i]=')' 且 s[i−1]=')'，也就是字符串形如 ".......))" ，我们可以推出：
如果 s[i−dp[i−1]−1]='(' ，那么dp[i]=dp[i−1]+dp[i−dp[i−1]−2]+2

这背后的原因是如果倒数第二个 ')' 是一个有效子字符串的一部分（记为 sub_s），对于最后一个 ')' ，
如果它是一个更长子字符串的一部分，那么它一定有一个对应的 '(' ，它的位置在倒数第二个 ')' 所在的有效子字符串的前面（也就是 sub_s的前面）。
因此，如果子字符串 sub_s的前面恰好是 '(' ，那么我们就用 22 加上 sub_s的长度（dp[i−1]）去更新 dp[i]。
除此以外，我们也会把有效子字符串 "(,sub_s,)"之前的有效子字符串的长度也加上，也就是加上 dp[i−dp[i−1]−2] 。

 */
fun longestValidParentheses(s: String): Int {
  var maxans = 0
  val dp = IntArray(s.length)
  for (i in 1 until s.length) {
    if (s[i] == ')') {
      if (s[i - 1] == '(') {
        dp[i] = (if (i >= 2) dp[i - 2] else 0) + 2
      } else if (i - dp[i - 1] > 0 && s[i - dp[i - 1] - 1] == '(') {
        dp[i] = dp[i - 1] + (if (i - dp[i - 1] >= 2) dp[i - dp[i - 1] - 2] else 0) + 2
      }
      maxans = kotlin.math.max(maxans, dp[i])
    }
  }
  return maxans
}
