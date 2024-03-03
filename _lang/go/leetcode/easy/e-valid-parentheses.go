package easy

/*
20. 有效的括号
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。

示例 1：

输入：s = "()"
输出：true
示例 2：

输入：s = "()[]{}"
输出：true
示例 3：

输入：s = "(]"
输出：false
示例 4：

输入：s = "([)]"
输出：false
示例 5：

输入：s = "{[]}"
输出：true

提示：

1 <= s.length <= 10^4
s 仅由括号 '()[]{}' 组成

链接：https://leetcode-cn.com/problems/valid-parentheses
*/
func isValid(s string) bool {
	l := len(s)
	if l&1 == 1 {
		return false
	}
	b1, b2, b3 := 0, 0, 0
	stack, si := make([]byte, l/2), 0
	skip := false
	for i := 0; i < l; i++ {
		c := s[i]
		switch c {
		case '(':
			b1++
		case ')':
			b1--
		case '[':
			b2++
		case ']':
			b2--
		case '{':
			b3++
		case '}':
			b3--
		}
		if b1 < 0 || b2 < 0 || b3 < 0 {
			return false
		}
		if i == 0 {
			stack[si] = s[i]
			continue
		}
		if skip {
			skip = false
			continue
		}
		if s[i] > stack[si] && s[i]-stack[si] < 3 {
			if si == 0 {
				if i+1 == l {
					return true
				} else {
					skip = true
					stack[si] = s[i+1]
				}
			} else {
				si--
			}
			continue
		}
		if i+1 == l {
			return false
		}
		si++
		if si == len(stack) {
			return false
		}
		stack[si] = s[i]
	}
	return false
}
