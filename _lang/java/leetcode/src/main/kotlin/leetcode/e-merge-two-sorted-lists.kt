package leetcode


/**
 *
21. 合并两个有序链表

将两个升序链表合并为一个新的升序链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。

示例：

输入：1->2->4, 1->3->4
输出：1->1->2->3->4->4

https://leetcode-cn.com/problems/merge-two-sorted-lists/
 */
//太简陋了这个连表，本来还想实现Iterator
class ListNode(var `val`: Int) {
  var next: ListNode? = null
  override fun toString(): String {
    var result = `val`.toString()
    var node = this
    while(node.next != null) {
      node = node.next!!
      result +="-> ${node.`val`}"
    }
    return result
  }
}


fun mergeTwoLists(l1: ListNode?, l2: ListNode?): ListNode? {
  if (l1 == null) return l2
  if (l2 == null) return l1
  var node1 = l1
  var node2 = l2
  var headNode: ListNode
  if (l1.`val` < l2.`val`) {
    headNode = ListNode(l1.`val`)
    node1 = l1.next
  } else {
    headNode = ListNode(l2.`val`)
    node2 = l2.next
  }
  val ans: ListNode = headNode
  while (true) {
    if (node1 == null) {
      headNode.next = node2
      return ans
    }
    if (node2 == null) {
      headNode.next = node1
      return ans
    }
    if (node1.`val` < node2.`val`) {
      headNode.next = ListNode(node1.`val`)
      headNode = headNode.next!!
      node1 = node1.next
    } else {
      headNode.next = ListNode(node2.`val`)
      headNode = headNode.next!!
      node2 = node2.next
    }
  }
}

fun mergeTwoListsV2(l1: ListNode?, l2: ListNode?): ListNode? {
  if (l1 == null) return l2
  if (l2 == null) return l1
  var node1 = l1
  var node2 = l2
  var headNode = ListNode(0)
  val ans = headNode
  while(node1!=null && node2 !=null){
    if (node1.`val` < node2.`val`) {
      headNode.next = node1
      node1 = node1.next
    } else {
      headNode.next = node2
      node2 = node2.next
    }
    headNode = headNode.next!!
  }
  if (node1 == null) headNode.next = node2 else headNode.next = node1
  return ans.next
}
