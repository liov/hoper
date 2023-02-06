package leetcode

/**
1014. 最佳观光组合

给定正整数数组 A，A[[i]] 表示第 i 个观光景点的评分，并且两个景点 i 和 j 之间的距离为 j - i。

一对景点（i < j）组成的观光组合的得分为（A[[i]] + A[[j]] + i - j）：景点的评分之和减去它们两者之间的距离。

返回一对观光景点能取得的最高分。



示例：

输入：[8,1,5,2,6]
输出：11
解释：i = 0, j = 2, A[[i]] + A[[j]] + i - j = 8 + 5 + 0 - 2 = 11


提示：

2 <= A.length <= 50000
1 <= A[[i]] <= 1000

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/best-sightseeing-pair
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。

在此题目中，设 res 为最终结果，dp[[i]] 为位置 i 对应的最高分

则有 dp[[0]] = 0, dp[[1]] = A[[1]] + A[[0]] - 1;

对于 i >= 2, dp[[i]] 仅取决于与 A[i - 1] 的分数 和 相对于 dp[i - 1] 的景点的分数，即

dp[[i]] = max(dp[[i - 1]] - A[[i - 1]] + A[[i]] - 1, A[[i]] + A[[i - 1]]- 1);

由于动态方程只用到了 dp[[i - 1]], 则可以只使用一个变量 last 来代替即可，实现 O(1) 空间

math{
j !=i-1 =》dp[[i-1]]+(A[[i]]-A[[i-1]]-1)
j == i-1 =》A[[i]]+A[[i-i]]-1
}

last = kotlin.math.max(last + (A[[i]] - A[[i - 1]] - 1), A[[i]] + A[[i - 1]] - 1) =
kotlin.math.max(last  - A[[i - 1]],  A[[i - 1]] ) + A[[i]] - 1=
if (last > 2 * A[[i - 1]]) last - A[[i - 1]] + A[[i]] - 1 else A[[i - 1]] + A[[i]] - 1
 */

fun maxScoreSightseeingPair(A: IntArray): Int {
  var last = A[0] + A[1] - 1
  var max = last
  for (i in 2 until A.size) {
    last = kotlin.math.max(last + (A[i] - A[i - 1] - 1), A[i] + A[i - 1] - 1)
    max = kotlin.math.max(max, last)
  }
  return max
}
