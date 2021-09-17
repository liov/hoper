package leetcode

/**
22. 括号生成

数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且 有效的 括号组合。



示例：

输入：n = 3
输出：[
"((()))",
"(()())",
"(())()",
"()(())",
"()()()"
]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/generate-parentheses
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun generateParenthesis(n: Int): List<String> {
  val ans = ArrayList<String>()
  val arr = CharArray(2 * n)
  add(arr,0,ans,0)
  return ans
}

fun add(s: CharArray, pos: Int, ans: ArrayList<String>, _balance: Int) {
  if (pos == s.size ) {
    if (_balance == 0) ans.add(String(s)) else return
  }
  else {
    s[pos] = '('
    var balance = _balance + 1 //开销啊，为什么参数只能不可变啊
    add(s, pos + 1, ans, balance)
    s[pos] = ')'
    balance = _balance - 1
    if (balance < 0) return
    add(s, pos + 1, ans, balance)
  }
}
