package leetcode

import java.util.*
import kotlin.collections.ArrayList

class TreeNode(var `val`: Int) {
  var left: TreeNode? = null
  var right: TreeNode? = null
  override fun toString(): String {
    return "node($left<- $`val` ->$right)"
  }
}

/**
二叉树的序列化与反序列化
序列化是将一个数据结构或者对象转换为连续的比特位的操作，进而可以将转换后的数据存储在一个文件或者内存中，同时也可以通过网络传输到另一个计算机环境，采取相反方式重构得到原数据。

请设计一个算法来实现二叉树的序列化与反序列化。这里不限定你的序列 / 反序列化算法执行逻辑，你只需要保证一个二叉树可以被序列化为一个字符串并且将这个字符串反序列化为原始的树结构。

示例:

你可以将以下二叉树：

1
/ \
2   3
/ \
4   5

序列化为 "[1,2,3,null,null,4,5]"
提示: 这与 LeetCode 目前使用的方式一致，详情请参阅 LeetCode 序列化二叉树的格式。你并非必须采取这种方式，你也可以采用其他的方法解决这个问题。

说明: 不要使用类的成员 / 全局 / 静态变量来存储状态，你的序列化和反序列化算法应该是无状态的。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/serialize-and-deserialize-binary-tree
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
//满二叉树和普通二叉树的数组形式不一样！！！
// Encodes a URL to a shortened URL.
fun serialize(root: TreeNode?): String {
  if (root == null) return "null"
  val ret = ArrayList<TreeNode?>()
  ret.add(root)
  var m = 0
  var n = 0
  var goon = false
  while (true) {
    for (j in m..n) {
      if (ret[j] == null) continue
      ret.add(ret[j]!!.left)
      ret.add(ret[j]!!.right)
      if (!goon && (ret[ret.size - 1] != null || ret[ret.size - 2] != null)) goon = true
    }
    m = n + 1
    n = ret.size - 1
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
      deque.first().left = TreeNode(arr[i + 1].toInt())
      deque.addLast(deque.first().left!!)
    }
    if (i + 2 < arr.size && arr[i + 2] != "null") {
      deque.first().right = TreeNode(arr[i + 2].toInt())
      deque.addLast(deque.first().right!!)
    }
    deque.removeFirst()
    i += 2
  }

  return root
}

