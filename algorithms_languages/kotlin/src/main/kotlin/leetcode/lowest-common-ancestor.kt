package leetcode

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
