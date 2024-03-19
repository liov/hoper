package middle

/*
22. 括号生成

中等

数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且 有效的 括号组合。



示例 1：

输入：n = 3
输出：["((()))","(()())","(())()","()(())","()()()"]
示例 2：

输入：n = 1
输出：["()"]


提示：

1 <= n <= 8
*/

func generateParenthesis(n int) []string {
	n *= 2
	var ans []string
	tmp := make([]byte, 0, n)
	var dfs func(i, balance int)
	dfs = func(i, balance int) {
		if balance < 0 {
			return
		}
		if i == n {
			if balance == 0 {
				ans = append(ans, string(tmp))
			}
			return
		}
		tmp = append(tmp, '(')
		dfs(i+1, balance+1)
		tmp = tmp[:len(tmp)-1]
		tmp = append(tmp, ')')
		dfs(i+1, balance-1)
		tmp = tmp[:len(tmp)-1]
	}
	dfs(0, 0)
	return ans
}
