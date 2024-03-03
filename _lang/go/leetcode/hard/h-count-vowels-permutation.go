package hard

/*
1220. 统计元音字母序列的数目
给你一个整数 n，请你帮忙统计一下我们可以按下述规则形成多少个长度为 n 的字符串：

字符串中的每个字符都应当是小写元音字母（'a', 'e', 'i', 'o', 'u'）
每个元音 'a' 后面都只能跟着 'e'
每个元音 'e' 后面只能跟着 'a' 或者是 'i'
每个元音 'i' 后面 不能 再跟着另一个 'i'
每个元音 'o' 后面只能跟着 'i' 或者是 'u'
每个元音 'u' 后面只能跟着 'a'
由于答案可能会很大，所以请你返回 模 10^9 + 7 之后的结果。



示例 1：

输入：n = 1
输出：5
解释：所有可能的字符串分别是："a", "e", "i" , "o" 和 "u"。
示例 2：

输入：n = 2
输出：10
解释：所有可能的字符串分别是："ae", "ea", "ei", "ia", "ie", "io", "iu", "oi", "ou" 和 "ua"。
示例 3：

输入：n = 5
输出：68


提示：

1 <= n <= 2 * 10^4

https://leetcode-cn.com/problems/count-vowels-permutation/
*/
/*
我们设 dp[i][j] 代表当前长度为 i 且以字符 j 为结尾的字符串的数目，其中在此 j=0,1,2,3,4 分别代表元音字母 ‘a’,‘e’,‘i’,‘o’,‘u’，
通过以上的字符规则，我们可以得到动态规划递推公式如下

dp[i][0]=dp[i−1][1]+dp[i−1][2]+dp[i−1][4]
dp[i][1]=dp[i−1][0]+dp[i−1][2]
dp[i][2]=dp[i−1][1]+dp[i−1][3]
dp[i][3]=dp[i−1][2]
dp[i][4]=dp[i−1][2]+dp[i−1][3]
∑i=0-4 dp[n][i]
*/
func countVowelPermutation(n int) int {
	const mod int = 1e9 + 7
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, 5)
	}
	for i := 0; i < 5; i++ {
		dp[0][i] = 1
	}
	for i := 1; i < n; i++ {
		dp[i][0] = dp[i-1][1] + dp[i-1][2] + dp[i-1][4]
		dp[i][0] %= mod
		dp[i][1] = dp[i-1][0] + dp[i-1][2]
		dp[i][1] %= mod
		dp[i][2] = dp[i-1][1] + dp[i-1][3]
		dp[i][2] %= mod
		dp[i][3] = dp[i-1][2]
		dp[i][3] %= mod
		dp[i][4] = dp[i-1][2] + dp[i-1][3]
		dp[i][4] %= mod
	}
	return (dp[n-1][0] + dp[n-1][1] + dp[n-1][2] + dp[n-1][3] + dp[n-1][4]) % mod
}
