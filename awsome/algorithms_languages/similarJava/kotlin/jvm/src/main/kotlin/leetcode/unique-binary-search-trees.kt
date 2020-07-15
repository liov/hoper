package leetcode

/**
 给定一个整数 n，求以 1 ... n 为节点组成的二叉搜索树有多少种？

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
