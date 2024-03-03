package middle

/*
5. 最长回文子串

给你一个字符串 s，找到 s 中最长的回文子串。

示例 1：

输入：s = "babad"
输出："bab"
解释："aba" 同样是符合题意的答案。
示例 2：

输入：s = "cbbd"
输出："bb"
示例 3：

输入：s = "a"
输出："a"
示例 4：

输入：s = "ac"
输出："a"

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/longest-palindromic-substring
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
func longestPalindrome(s string) string {
	ret := s[0:1]
	for i := range s {
		j := 1

		for {
			if i-j+1 < 0 || i+j >= len(s) {
				break
			}
			if s[i-j+1] == s[i+j] {
				if j*2 > len(ret) {
					ret = s[i-j+1 : i+j+1]
				}
			} else {
				break
			}
			j++
		}
		j = 1
		for {
			if i-j < 0 || i+j >= len(s) {
				break
			}
			if s[i-j] == s[i+j] {
				if j*2+1 > len(ret) {
					ret = s[i-j : i+j+1]
				}
			} else {
				break
			}
			j++
		}
	}
	return ret
}

func longestPalindrome2(s string) string {
	ret := s[0:1]
	for i := range s {
		j := 1

		for {
			if i-j+1 < 0 || i+j >= len(s) {
				break
			}
			if s[i-j+1] == s[i+j] {
				if j*2 > len(ret) {
					ret = s[i-j+1 : i+j+1]
				}
			} else {
				break
			}
			j++
		}
		j = 1
		for {
			if i-j < 0 || i+j >= len(s) {
				break
			}
			if s[i-j] == s[i+j] {
				if j*2+1 > len(ret) {
					ret = s[i-j : i+j+1]
				}
			} else {
				break
			}
			j++
		}
	}
	return ret
}
