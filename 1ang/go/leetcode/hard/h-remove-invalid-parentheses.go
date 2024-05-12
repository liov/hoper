package hard

import (
	"math/bits"
	"test/leetcode"
)

/*
301. 删除无效的括号
给你一个由若干括号和字母组成的字符串 s ，删除最小数量的无效括号，使得输入的字符串有效。

返回所有可能的结果。答案可以按 任意顺序 返回。



示例 1：

输入：s = "()())()"
输出：["(())()","()()()"]
示例 2：

输入：s = "(a)())()"
输出：["(a())()","(a)()()"]
示例 3：

输入：s = ")("
输出：[""]


提示：

1 <= s.length <= 25
s 由小写英文字母以及括号 '(' 和 ')' 组成
s 中至多含 20 个括号

https://leetcode-cn.com/problems/remove-invalid-parentheses/
*/

func checkValid(str []byte, lmask, rmask int, left, right []int) bool {
	cnt := 0
	pos1, pos2 := 0, 0
	for i := range str {
		if pos1 < len(left) && i == left[pos1] {
			if lmask>>pos1&1 == 0 {
				cnt++
			}
			pos1++
		} else if pos2 < len(right) && i == right[pos2] {
			if rmask>>pos2&1 == 0 {
				cnt--
				if cnt < 0 {
					return false
				}
			}
			pos2++
		}
	}
	return cnt == 0
}

func recoverStr(str []byte, lmask, rmask int, left, right []int) string {
	var res []byte
	pos1, pos2 := 0, 0
	for i, ch := range str {
		if pos1 < len(left) && i == left[pos1] {
			if lmask>>pos1&1 == 0 {
				res = append(res, ch)
			}
			pos1++
		} else if pos2 < len(right) && i == right[pos2] {
			if rmask>>pos2&1 == 0 {
				res = append(res, ch)
			}
			pos2++
		} else {
			res = append(res, ch)
		}
	}
	return string(res)
}

func removeInvalidParentheses(s string) (ans []string) {
	byteses := []byte(s)
	buf1 := make([]byte, 0, len(s))
	firstClosingParenthesis := false
	for _, c := range byteses {
		if firstClosingParenthesis {
			buf1 = append(buf1, c)
			continue
		}
		if c >= 'a' && c <= 'z' {
			buf1 = append(buf1, c)
			continue
		}
		if c == '(' && !firstClosingParenthesis {
			firstClosingParenthesis = true
			buf1 = append(buf1, c)
		}
	}
	buf2 := make([]byte, 0, len(buf1))
	firstOpenParenthesis := false
	for i := len(buf1) - 1; i >= 0; i-- {
		if firstOpenParenthesis {
			buf2 = append(buf2, buf1[i])
			continue
		}
		if buf1[i] >= 'a' && buf1[i] <= 'z' {
			buf2 = append(buf2, buf1[i])
			continue
		}
		if buf1[i] == ')' && !firstOpenParenthesis {
			firstOpenParenthesis = true
			buf2 = append(buf2, buf1[i])
		}
	}
	byteses = leetcode.reverse(buf2)
	var left, right []int
	lremove, rremove := 0, 0
	for i, ch := range byteses {
		if ch == '(' {
			left = append(left, i)
			lremove++
		} else if ch == ')' {
			right = append(right, i)
			if lremove == 0 {
				rremove++
			} else {
				lremove--
			}
		}
	}

	var maskArr1, maskArr2 []int
	for i := 0; i < 1<<len(left); i++ {
		if bits.OnesCount(uint(i)) == lremove {
			maskArr1 = append(maskArr1, i)
		}
	}
	for i := 0; i < 1<<len(right); i++ {
		if bits.OnesCount(uint(i)) == rremove {
			maskArr2 = append(maskArr2, i)
		}
	}

	res := map[string]struct{}{}
	for _, mask1 := range maskArr1 {
		for _, mask2 := range maskArr2 {
			if checkValid(byteses, mask1, mask2, left, right) {
				res[recoverStr(byteses, mask1, mask2, left, right)] = struct{}{}
			}
		}
	}
	for str := range res {
		ans = append(ans, str)
	}
	return
}

func invalidParentheses(str string) bool {
	cnt := 0
	for _, ch := range str {
		if ch == '(' {
			cnt++
		} else if ch == ')' {
			cnt--
			if cnt < 0 {
				return false
			}
		}
	}
	return cnt == 0
}

func removeInvalidParentheses2(s string) (ans []string) {
	curSet := map[string]struct{}{s: {}}
	for {
		for str := range curSet {
			if invalidParentheses(str) {
				ans = append(ans, str)
			}
		}
		if len(ans) > 0 {
			return
		}
		nextSet := map[string]struct{}{}
		for str := range curSet {
			for i, ch := range str {
				if i > 0 && byte(ch) == str[i-1] {
					continue
				}
				if ch == '(' || ch == ')' {
					nextSet[str[:i]+str[i+1:]] = struct{}{}
				}
			}
		}
		curSet = nextSet
	}
}
