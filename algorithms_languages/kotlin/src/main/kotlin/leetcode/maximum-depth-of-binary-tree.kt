package leetcode

/**
 * 给定一个二叉树，找出其最大深度。
 */
fun maxDepth(root: TreeNode?): Int {
  return if (root == null)  0 else kotlin.math.max(maxDepth(root.left),maxDepth(root.right)) + 1
}
