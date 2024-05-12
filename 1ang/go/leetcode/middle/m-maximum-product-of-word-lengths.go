package middle

import "sort"

/*
318. 最大单词长度乘积
给定一个字符串数组 words，找到 length(word[i]) * length(word[j]) 的最大值，并且这两个单词不含有公共字母。你可以认为每个单词只包含小写字母。如果不存在这样的两个单词，返回 0。

示例 1:

输入: ["abcw","baz","foo","bar","xtfn","abcdef"]
输出: 16
解释: 这两个单词为 "abcw", "xtfn"。
示例 2:

输入: ["a","ab","abc","d","cd","bcd","abcd"]
输出: 4
解释: 这两个单词为 "ab", "cd"。
示例 3:

输入: ["a","aa","aaa","aaaa"]
输出: 0
解释: 不存在这样的两个单词。

提示：

2 <= words.length <= 1000
1 <= words[i].length <= 1000
words[i] 仅包含小写字母

https://leetcode-cn.com/problems/maximum-product-of-word-lengths/
*/
func maxProduct(words []string) int {
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) > len(words[j])
	})
	var ans int
	for i := 0; i < len(words)-1; i++ {
		if ans > len(words[i])*len(words[i+1]) {
			return ans
		}
		arr := make([]bool, 26)
		for _, c := range words[i] {
			arr[c-'a'] = true
		}
		for j := i + 1; j < len(words); j++ {
			if ans > len(words[i])*len(words[j]) {
				continue
			}
			eff := true
			for _, c := range words[j] {
				if arr[c-'a'] {
					eff = false
					break
				}
			}
			if eff {
				ans = max(ans, len(words[i])*len(words[j]))
			}
		}
	}
	return ans
}

func maxProduct2(words []string) int {
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) > len(words[j])
	})
	var ans int
	for i := 0; i < len(words)-1; i++ {
		if ans > len(words[i])*len(words[i+1]) {
			return ans
		}
		var arr int
		for _, c := range words[i] {
			arr |= 1 << (c - 'a')
		}
		for j := i + 1; j < len(words); j++ {
			if ans > len(words[i])*len(words[j]) {
				continue
			}
			eff := true
			for _, c := range words[j] {
				if arr>>(c-'a')&1 == 1 {
					eff = false
					break
				}
			}
			if eff {
				ans = max(ans, len(words[i])*len(words[j]))
			}
		}
	}
	return ans
}
