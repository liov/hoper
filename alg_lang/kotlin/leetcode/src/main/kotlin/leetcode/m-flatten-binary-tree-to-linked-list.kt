package leetcode

/**
114. 二叉树展开为链表
给你二叉树的根结点 root ，请你将它展开为一个单链表：

展开后的单链表应该同样使用 TreeNode ，其中 right 子指针指向链表中下一个结点，而左子指针始终为 null 。
展开后的单链表应该与二叉树 先序遍历 顺序相同。


示例 1：


输入：root = [1,2,5,3,4,null,6]
输出：[1,null,2,null,3,null,4,null,5,null,6]
示例 2：

输入：root = []
输出：[]
示例 3：

输入：root = [0]
输出：[0]


提示：

树中结点数在范围 [0, 2000] 内
-100 <= Node.val <= 100


进阶：你可以使用原地算法（O(1) 额外空间）展开这棵树吗？


 */
//给定一个二叉树，原地将它展开为一个单链表。
fun flatten(root: TreeNode?): Unit {
  var root = root
  while (root != null) {
    //左子树为 null，直接考虑下一个节点
    if (root.left == null)  root = root.right
    else {
      // 找左子树最右边的节点
      var pre = root.left!!
      while (pre.right != null) pre = pre.right!!
      //将原来的右子树接到左子树的最右边节点
      pre.right = root.right
      // 将左子树插入到右子树的地方
      root.right = root.left
      root.left = null
      // 考虑下一个节点
      root = root.right
    }
  }
}
