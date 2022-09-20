package leetcode

/**
392. 判断子序列

给定字符串 s 和 t ，判断 s 是否为 t 的子序列。

你可以认为 s 和 t 中仅包含英文小写字母。字符串 t 可能会很长（长度 ~= 500,000），而 s 是个短字符串（长度 <=100）。

字符串的一个子序列是原始字符串删除一些（也可以不删除）字符而不改变剩余字符相对位置形成的新字符串。（例如，"ace"是"abcde"的一个子序列，而"aec"不是）。

示例 1:
s = "abc", t = "ahbgdc"

返回 true.

示例 2:
s = "axc", t = "ahbgdc"

返回 false.

后续挑战 :

如果有大量输入的 S，称作S1, S2, ... , Sk 其中 k >= 10亿，你需要依次检查它们是否为 T 的子序列。在这种情况下，你会怎样改变代码？
 */
fun isSubsequence(s: String, t: String): Boolean {
  if (s == "") return true
  if (t.length < s.length) return false
  var i = 0
  for (c in t) {
    if (c == s[i]) i++
    if (i == s.length) return true
  }
  return false
}

fun isSubsequenceV2(s: String, t: String): Boolean {
  if (s == "") return true
  if (t.length < s.length) return false
  var i = 0
  var j = s.length - 1
  val even = t.length and 1 == 0

  for (n in 0 until t.length / 2) {
    if (t[n] == s[i]) {
      if (i == j) return true
      i++
    }
    if (t[t.length - 1 - n] == s[j]) {
      if (i == j) return true
      j--
    }
  }
  if (!even) if (t[t.length / 2] == s[i]) if (i == j) return true
  return false
}
