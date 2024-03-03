package middle

import "strings"

/*
686. 重复叠加字符串匹配
给定两个字符串 a 和 b，寻找重复叠加字符串 a 的最小次数，使得字符串 b 成为叠加后的字符串 a 的子串，如果不存在则返回 -1。

注意：字符串 "abc" 重复叠加 0 次是 ""，重复叠加 1 次是 "abc"，重复叠加 2 次是 "abcabc"。

示例 1：

输入：a = "abcd", b = "cdabcdab"
输出：3
解释：a 重复叠加三遍后为 "abcdabcdabcd", 此时 b 是其子串。
示例 2：

输入：a = "a", b = "aa"
输出：2
示例 3：

输入：a = "a", b = "a"
输出：1
示例 4：

输入：a = "abc", b = "wxyz"
输出：-1

提示：

1 <= a.length <= 10^4
1 <= b.length <= 10^4
a 和 b 由小写英文字母组成

https://leetcode-cn.com/problems/repeated-string-match/
*/
func repeatedStringMatch(a string, b string) int {
	if b == "" {
		return 0
	}
	if a == "" {
		return -1
	}
	if a == b || (len(a) > len(b) && strings.Contains(a, b)) {
		return 1
	}

	var bitmap int
	for i := range a {
		bitmap |= 1 << (a[i] - 'a')
	}
	start, end, ans := -1, -1, 0
	for i := 0; i < len(b); i++ {
		if bitmap>>(b[i]-'a')&1 != 1 {
			return -1
		}
		if b[i] == a[0] {
			if strings.HasPrefix(a, b[i:]) && strings.HasSuffix(a, b[:i]) {
				return 2
			}
			if i+len(a) <= len(b) {
				if b[i:i+len(a)] == a {
					ans++
					end = i + len(a)
					if start == -1 {
						start = i
						if start > 0 {
							if strings.HasSuffix(a, b[:start]) {
								ans++
							} else {
								return -1
							}
						}
					}
					i = end - 1
				} else {
					if end != -1 {
						return -1
					}
				}
			}
		}
	}
	if end < len(b) {
		if start != -1 && strings.HasPrefix(a, b[end:]) {
			ans++
		} else {
			return -1
		}
	}

	return ans
}
