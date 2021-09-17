package leetcode

/**
108. 将有序数组转换为二叉搜索树

将一个按照升序排列的有序数组，转换为一棵高度平衡二叉搜索树。

本题中，一个高度平衡二叉树是指一个二叉树每个节点 的左右两个子树的高度差的绝对值不超过 1。

示例:

给定有序数组: [-10,-3,0,5,9],

一个可能的答案是：[0,-3,9,-10,null,5]，它可以表示下面这个高度平衡二叉搜索树：

0
/ \
-3   9
/   /
-10  5

https://leetcode-cn.com/problems/convert-sorted-array-to-binary-search-tree/
 */
fun sortedArrayToBST(nums: IntArray): TreeNode? {
  return helper(nums, 0, nums.size - 1)
}

fun helper(nums: IntArray, left: Int, right: Int): TreeNode? {
  if (left > right) return null
  val mid = (left + right + 1) / 2
  val root = TreeNode(nums[mid])
  if (left == right) return root
  root.left = helper(nums, left, mid - 1)
  root.right = helper(nums, mid + 1, right)
  return root
}
