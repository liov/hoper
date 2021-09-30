package leetcode

/**
不同的二叉搜索树

给定一个整数 n，求以 1 ... n 为节点组成的二叉搜索树有多少种？

示例:

输入: 3
输出: 5
解释:
给定 n = 3, 一共有 5 种不同结构的二叉搜索树:

1         3     3      2      1
\       /     /      / \      \
3     2     1      1   3      2
/     /       \                 \
2     1         2                 3

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/unique-binary-search-trees
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
//Gn = sum(i in 1..n){G(i-1)G(n-i)}
//Gn+1=(2(2n+1)/(n+2))Gn
fun numTrees(n: Int): Int {
  var c = 1
  for (i in 0 until n) {
    c = c * 2 * (2 * i + 1) / (i + 2)
  }
  return c
}

fun numTreesV2(n: Int): Int {
  val arr = IntArray(n + 1)
  arr[0] = 1
  arr[1] = 1
  for (i in 2..n) {
    for (j in 1..i) arr[i] += arr[j - 1] * arr[i - j]
  }
  return arr[n]
}

/**
 * 生成二叉搜索树
 */
fun generateTrees(n: Int): List<TreeNode?> {
  if (n == 1) return listOf(TreeNode(n))
  return helper(1, n)
}

fun helper(start: Int, end: Int): List<TreeNode?> {
  if (start > end) return listOf(null)
  if (start == end) return listOf(TreeNode(start))
  val list = ArrayList<TreeNode?>()
  for (i in start..end) {
    val left = helper(start, i - 1)
    val right = helper(i + 1, end)
    for (l in left.indices) {
      for (r in right.indices) {
        val tree = TreeNode(i)
        tree.left = left[l]
        tree.right = right[r]
        list.add(tree)
      }
    }
  }
  return list
}
