package leetcode

import java.util.*


/**
23. 合并K个升序链表

给你一个链表数组，每个链表都已经按升序排列。

请你将所有链表合并到一个升序链表中，返回合并后的链表。



示例 1：

输入：lists = [[1,4,5],[1,3,4],[2,6]]
输出：[1,1,2,3,4,4,5,6]
解释：链表数组如下：
[
1->4->5,
1->3->4,
2->6
]
将它们合并到一个有序链表中得到。
1->1->2->3->4->4->5->6
示例 2：

输入：lists = []
输出：[]
示例 3：

输入：lists = [[]]
输出：[]


提示：

k == lists.length
0 <= k <= 10^4
0 <= lists[i].length <= 500
-10^4 <= lists[i][j] <= 10^4
lists[i] 按 升序 排列
lists[i].length 的总和不超过 10^4

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/merge-k-sorted-lists
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
//if 要加括号就很不爽
fun mergeKLists(lists: Array<ListNode?>): ListNode? {
  if (lists.isEmpty()) return null
  if (lists.size == 1) return lists[0]
  val tree = BinaryTree<Int>()
  var tmp: ListNode?
  for (i in lists.indices) {
    tmp = lists[i]
    while (tmp != null) {
      tree.insert(tmp.`val`)
      tmp = tmp.next
    }
  }
  val ans = ListNode(0)
  tmp = ans
  tree.sequence().forEach {
    tmp!!.next = ListNode(it)
    tmp= tmp!!.next
  }
  return ans.next
}

class BinaryTree<T : Comparable<T>> {
  private var root: Node<T>? = null
  var size = 0

  constructor()
  constructor(v: T) {
    root = Node(v)
    size = 1
  }

  fun isEmpty(): Boolean {
    return size == 0
  }

  fun insert(v: T) {
    size += 1
    if (root == null) {
      root = Node(v)
      return
    }
    root!!.insert(v)
  }

  fun insertLoop(v: T) {
    size += 1
    if (root == null) {
      root = Node(v)
      return
    }
    var tmp = root
    while (true) {
      if (v > tmp!!.value) {
        if (tmp.right == null) {
          tmp.right = Node(v)
          return
        } else tmp = tmp.right
      } else {
        if (tmp.left == null) {
          tmp.left = Node(v)
          return
        } else tmp = tmp.left
      }
    }
  }


  /**
   * 前序遍历
   * @param node 开始遍历元素
   */
  fun prevRecursive(node: Node<T>? = root) {
    if (node != null) {
      print("${node.value} ")
      prevRecursive(node.left)
      prevRecursive(node.right)
    }
  }

  fun prevIterator(node: Node<T>? = root): List<T> {
    val result = ArrayList<T>()
    val stack = Stack<Node<T>>()
    if (root != null) stack.push(root) else return result
    var tmp: Node<T>?
    while (!stack.isEmpty()) {
      tmp = stack.pop()
      result.add(tmp.value)
      if (tmp.right != null) stack.push(tmp.right); // 右节点入栈
      if (tmp.left != null) stack.push(tmp.left); // 左节点入栈
    }
    return result
  }

  /**
   * 中序遍历
   * @param node 开始遍历元素
   */
  fun midRecursive(node: Node<T>? = root) {
    if (node != null) {
      midRecursive(node.left)
      print("${node.value} ")
      midRecursive(node.right)
    }
  }

  fun midIterator(): List<T> {
    val result = ArrayList<T>()
    val stack = Stack<Node<T>>()
    var tmp = root
    while (!stack.isEmpty() || tmp != null) {
      while (tmp != null) {
        stack.push(tmp); // 添加根节点
        tmp = tmp.left; // 循环添加左节点
      }
      tmp = stack.pop()
      result.add(tmp.value)
      tmp = tmp.right
    }
    return result
  }

  fun sequence(): Sequence<T> = sequence {
    val stack = Stack<Node<T>>()
    var tmp = root
    while (!stack.isEmpty() || tmp != null) {
      while (tmp != null) {
        stack.push(tmp); // 添加根节点
        tmp = tmp.left; // 循环添加左节点
      }
      tmp = stack.pop()
      yield(tmp.value)
      tmp = tmp.right
    }
  }

  /**
   * 后序遍历
   * @param node 开始遍历元素
   */
  fun subRecursive(node: Node<T>? = root) {
    if (node != null) {
      subRecursive(node.left);
      subRecursive(node.right);
      print("${node.value} ");
    }
  }
}

class Node<T : Comparable<T>>(val value: T) {
  var left: Node<T>? = null
  var right: Node<T>? = null
  fun insert(v: T) {
    if (v > this.value) {
      if (this.right == null) this.right = Node(v) else this.right!!.insert(v)
    } else {
      if (this.left == null) this.left = Node(v) else this.left!!.insert(v)
    }
  }
}

fun mergeKListsV2(lists: Array<ListNode?>): ListNode? {
  return merge(lists, 0, lists.size - 1);
}

fun merge(lists:Array<ListNode?>,l:Int,r:Int): ListNode? {
  if (l == r) return lists[l];
  if (l > r) return null
  val mid = (l + r) shr 1
  return mergeTwoListsV2(merge(lists, l, mid), merge(lists, mid + 1, r))
}
