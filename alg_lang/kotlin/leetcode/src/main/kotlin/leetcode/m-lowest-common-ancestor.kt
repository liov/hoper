package leetcode

/*
236. 二叉树的最近公共祖先
给定一个二叉树, 找到该树中两个指定节点的最近公共祖先。

百度百科中最近公共祖先的定义为：“对于有根树 T 的两个节点 p、q，最近公共祖先表示为一个节点 x，满足 x 是 p、q 的祖先且 x 的深度尽可能大（一个节点也可以是它自己的祖先）。”



示例 1：


输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 1
输出：3
解释：节点 5 和节点 1 的最近公共祖先是节点 3 。
示例 2：


输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 4
输出：5
解释：节点 5 和节点 4 的最近公共祖先是节点 5 。因为根据定义最近公共祖先节点可以为节点本身。
示例 3：

输入：root = [1,2], p = 1, q = 2
输出：1


提示：

树中节点数目在范围 [2, 105] 内。
-109 <= Node.val <= 109
所有 Node.val 互不相同 。
p != q
p 和 q 均存在于给定的二叉树中。

https://leetcode-cn.com/problems/lowest-common-ancestor-of-a-binary-tree/
 */
fun lowestCommonAncestor(root: TreeNode?, node1: TreeNode, node2: TreeNode): TreeNode? {
  if (root == null || root == node1 || root == node2)  return root
  //采用递归调用的思路，将二叉树分为左子树和右子树分别处理
  //使用递归需要注意两点：
  //1.子问题需与原始问题为同样的问题，且更为简单；2.不能无限制地调用本身，程序必须有出口
  //查看左子树中是否有目标结点，没有为null
  val left = lowestCommonAncestor(root.left, node1, node2)
  //查看右子树是否有目标节点，没有为null
  val right = lowestCommonAncestor(root.right, node1, node2)
  //都不为空，说明做右子树都有目标结点，则公共祖先就是本身
  return if (left == node1 && right == node2 || left == node2 && right == node1) root else left ?: right
}
