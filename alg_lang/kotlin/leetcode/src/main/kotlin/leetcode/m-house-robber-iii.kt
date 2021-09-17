package leetcode

/**
337. 打家劫舍 III

这个地区只有一个入口，我们称之为“根”。 除了“根”之外，每栋房子有且只有一个“父“房子与之相连。一番侦察之后，聪明的小偷意识到“这个地方的所有房屋的排列类似于一棵二叉树”。 如果两个直接相连的房子在同一天晚上被打劫，房屋将自动报警。

计算在不触动警报的情况下，小偷一晚能够盗取的最高金额。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/house-robber-iii
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun rob(root: TreeNode?): Int {
  if(root == null) return 0
  val sub1 = root.`val` + robNotContain(root.left)+robNotContain(root.right)
  val sub2 = rob(root.left)+rob(root.right)
  return kotlin.math.max(sub1,sub2)
}

fun robNotContain(root: TreeNode?): Int {
  if(root == null) return 0
  return rob(root.left) + rob(root.right)
}

fun robV2(root: TreeNode?): Int {
  val rootStatus = dfs(root)
  return kotlin.math.max(rootStatus[0], rootStatus[1])
}

fun dfs(node:TreeNode?):IntArray {
  if (node == null) return IntArray(2)
  val l = dfs(node.left)
  val r = dfs(node.right)
  val selected = node.`val` + l[1] + r[1]
  val notSelected = kotlin.math.max(l[0], l[1]) + kotlin.math.max(r[0], r[1]);
  return intArrayOf(selected, notSelected)
}
