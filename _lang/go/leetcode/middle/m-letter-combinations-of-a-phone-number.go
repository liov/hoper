package middle

/*
17. 电话号码的字母组合

中等

给定一个仅包含数字 2-9 的字符串，返回所有它能表示的字母组合。答案可以按 任意顺序 返回。

给出数字到字母的映射如下（与电话按键相同）。注意 1 不对应任何字母。





示例 1：

输入：digits = "23"
输出：["ad","ae","af","bd","be","bf","cd","ce","cf"]
示例 2：

输入：digits = ""
输出：[]
示例 3：

输入：digits = "2"
输出：["a","b","c"]


提示：

0 <= digits.length <= 4
digits[i] 是范围 ['2', '9'] 的一个数字。
*/

var numberWords = [][]byte{
	{}, {},
	{'a', 'b', 'c'},
	{'d', 'e', 'f'},
	{'g', 'h', 'i'},
	{'j', 'k', 'l'},
	{'m', 'n', 'o'},
	{'p', 'q', 'r', 's'},
	{'t', 'u', 'v'},
	{'w', 'x', 'y', 'z'},
}

func letterCombinations(digits string) []string {
	if digits == "" {
		return nil
	}
	n := len(digits)
	var ans []string
	tmp := make([]byte, 0)
	var dsf func(i int)
	dsf = func(i int) {
		if i == n {
			ans = append(ans, string(tmp))
			return
		}
		for j := 0; j < len(numberWords[digits[i]-'0']); j++ {
			c := numberWords[digits[i]-'0'][j]
			tmp = append(tmp, c)
			dsf(i + 1)
			tmp = tmp[:len(tmp)-1]
		}
	}
	dsf(0)
	return ans
}
