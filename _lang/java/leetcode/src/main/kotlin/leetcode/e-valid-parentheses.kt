package leetcode

/**
有效的括号

给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
注意空字符串可被认为是有效字符串。

示例 1:

输入: "()"
输出: true
示例 2:

输入: "()[]{}"
输出: true
示例 3:

输入: "(]"
输出: false
示例 4:

输入: "([)]"
输出: false
示例 5:

输入: "{[]}"
输出: true

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/valid-parentheses
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */

class LinkedListNode<T>(val value: T) {
  var next: LinkedListNode<T>? = null
  var prev: LinkedListNode<T>? = null
}

infix fun Char.less3(to: Char): Boolean {
  return this - to in 1..2
}

//本人的链表大法
//176 ms , 在所有 Kotlin 提交中击败了 72.22% 的用户
fun isValid(s: String): Boolean {
  if (s.length and 1 == 1) return false
  if (s == "" || (s.length == 2 && s[1] less3 s[0])) return true //ascii码表，差不超过2

  var currentNode = LinkedListNode(s[0])
  var skip = false
  for (i in 1 until s.length) {
    if (skip) {
      skip = false
      continue
    }
    if (s[i] less3 currentNode.value) {
      if (currentNode.prev == null) {
        if (i + 1 == s.length) return true
        currentNode = LinkedListNode(s[i + 1])
        skip = true
      } else currentNode = currentNode.prev!!
      continue
    }
    if (i + 1 == s.length) return false
    currentNode.next = LinkedListNode(s[i])
    currentNode.next!!.prev = currentNode
    currentNode = currentNode.next!!
  }
  return false
}

fun isValidV2(s: String): Boolean {
  if (s.length and 1 == 1) return false
  if (s == "" || (s.length == 2 && s[1] less3 s[0])) return true //ascii码表，差不超过2
  val stack = CharArray(s.length / 2)
  var skip = false
  var idx = 0 //手动标志位,go 来做是不是更爽
  stack[idx] = s[0]
  for (i in 1 until s.length) {
    if (skip) {
      skip = false
      continue
    }
    if (s[i] less3 stack[idx]) {
      if (idx == 0) {
        if (i + 1 == s.length) return true
        stack[idx] = s[i + 1]
        skip = true
      } else idx -= 1
      continue
    }
    if (i + 1 == s.length) return false
    idx += 1
    if (idx == stack.size) return false
    stack[idx] = s[i]
  }
  return false
}
