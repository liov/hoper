package leetcode

/*
524. 通过删除字母匹配到字典里最长单词
给你一个字符串 s 和一个字符串数组 dictionary 作为字典，找出并返回字典中最长的字符串，该字符串可以通过删除 s 中的某些字符得到。

如果答案不止一个，返回长度最长且字典序最小的字符串。如果答案不存在，则返回空字符串。



示例 1：

输入：s = "abpcplea", dictionary = ["ale","apple","monkey","plea"]
输出："apple"
示例 2：

输入：s = "abpcplea", dictionary = ["a","b","c"]
输出："a"


提示：

1 <= s.length <= 1000
1 <= dictionary.length <= 1000
1 <= dictionary[i].length <= 1000
s 和 dictionary[i] 仅由小写英文字母组成

https://leetcode-cn.com/problems/longest-word-in-dictionary-through-deleting

*/

func findLongestWord(s string, dictionary []string) string {
	var ret string
	for _, v := range dictionary {
		if len(v) >= len(ret) && findLongestWordHelper(s, v) {
			if len(v) > len(ret) {
				ret = v
			} else {
				for j := range v {
					if v[j] < ret[j] {
						ret = v
						break
					}
					if v[j] > ret[j] {
						break
					}
				}
			}
		}
	}
	return ret
}

func findLongestWordHelper(s string, dictionary string) bool {
	i := 0
	for j := 0; j < len(dictionary); j++ {
		for {
			if dictionary[j] == s[i] {
				if j == len(dictionary)-1 {
					return true
				}
				i++
				if i == len(s) {
					return false
				}
				break
			}
			i++
			if i == len(s) {
				return false
			}
		}
	}
	return false
}
