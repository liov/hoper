package leetcode

import java.util.*
import kotlin.collections.ArrayList

class TreeNode(var `val`: Int) {
  var left: TreeNode? = null
  var right: TreeNode? = null
}

//满二叉树和普通二叉树的数组形式不一样！！！
// Encodes a URL to a shortened URL.
fun serialize(root: TreeNode?): String {
  if (root == null) return "null"
  val ret = ArrayList<TreeNode?>()
  ret.add(root)
  var m = 0
  var n = 0
  var count = 0
  var goon = false
  while (true) {
    for (j in m..n) {
      if (ret[j] == null) continue
      ret.add(ret[j]!!.left)
      ret.add(ret[j]!!.right)
      count += 1
      if (!goon && (ret[ret.size - 1] != null || ret[ret.size - 2] != null)) goon = true
    }
    m = n + 1
    n = m + count
    count = 0
    if (goon) goon = false else break
  }
  //可以优化
  for (j in ret.size - 1 downTo 0) {
    if (ret[j] == null) ret.removeAt(j) else break
  }
  val separator = ","
  return ret.joinToString(separator) { it?.`val`.toString() }
}

// Decodes your encoded data to tree.

fun deserialize(data: String): TreeNode? {
  if (data == "null") return null
  val arr = data.split(",")
  val deque = ArrayDeque<TreeNode>()
  val root = TreeNode(arr[0].toInt())
  deque.addLast(root)
  var i = 0
  while (deque.isNotEmpty()) {
    if (i + 1 < arr.size && arr[i + 1] != "null") {
      deque.first().left= TreeNode(arr[i + 1].toInt())
      deque.addLast(deque.first().left!!)
    }
     if (i + 2 < arr.size && arr[i + 2] != "null") {
       deque.first().right = TreeNode(arr[i + 2].toInt())
       deque.addLast(deque.first().right!!)
    }
    deque.removeFirst()
    i+=2
  }

  return root
}

