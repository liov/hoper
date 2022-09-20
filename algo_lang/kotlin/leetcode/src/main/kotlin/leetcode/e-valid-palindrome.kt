package leetcode

/**
验证回文串

给定一个字符串，验证它是否是回文串，只考虑字母和数字字符，可以忽略字母的大小写。

说明：本题中，我们将空字符串定义为有效的回文串。

示例 1:

输入: "A man, a plan, a canal: Panama"
输出: true
示例 2:

输入: "race a car"
输出: false

https://leetcode-cn.com/problems/valid-palindrome/
 */
//与有效括号一样,栈都不用，直接双指针
fun isPalindrome(s: String): Boolean {
  if(s == "") return true
  var start = 0
  var end = s.length - 1
  val isLetter = {a: Char -> a in 'a'..'z' || a in 'A'..'Z'}
  val eqLetter ={a:Char,b:Char-> a==b || (isLetter(a) && isLetter(b) && (s[start] -32 == s[end] || s[end] -32 == s[start]))}
  while(true){
    //比起这种写法s[start] !in 'a'..'z' && s[start] !in 'A'..'Z' && s[start] !in '0'..'9'，效率高
    while (!(isLetter(s[start]) || s[start] in '0'..'9')) {
      if (start == end) return true
      start++
    }
    while (!(isLetter(s[end]) || s[end] in '0'..'9')) end--
    if(start >= end) return true
    if (eqLetter(s[start],s[end])) {
      start++
      end--
    }else return false
  }
}
