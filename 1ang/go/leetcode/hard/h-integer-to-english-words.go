package hard

import (
	"strings"
)

/*
将非负整数 num 转换为其对应的英文表示。

示例 1：

输入：num = 123
输出："One Hundred Twenty Three"
示例 2：

输入：num = 12345
输出："Twelve Thousand Three Hundred Forty Five"
示例 3：

输入：num = 1234567
输出："One Million Two Hundred Thirty Four Thousand Five Hundred Sixty Seven"
示例 4：

输入：num = 1234567891
输出："One Billion Two Hundred Thirty Four Million Five Hundred Sixty Seven Thousand Eight Hundred Ninety One"

提示：

0 <= num <= 2^31 - 1

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/integer-to-english-words
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
var (
	low  = []string{"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine"}
	mid  = []string{"Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}
	high = []string{"", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}
)

func numberToWords(num int) string {
	if num == 0 {
		return "Zero"
	}
	var ret strings.Builder
	div := num / 10e8
	if div > 0 {
		ret.WriteString(low[div])
		ret.WriteString(" Billion")
	}
	num = int(num % 10e8)
	div = num / 10e5
	if div > 0 {
		if ret.Len() > 0 {
			ret.WriteString(" ")
		}
		ret.WriteString(numberToWordsHelper(div))
		ret.WriteString(" Million")
	}
	num = num % 10e5
	div = num / 10e2
	if div > 0 {
		if ret.Len() > 0 {
			ret.WriteString(" ")
		}
		ret.WriteString(numberToWordsHelper(div))
		ret.WriteString(" Thousand")
	}
	num = num % 10e2
	if num > 0 {
		if ret.Len() > 0 {
			ret.WriteString(" ")
		}
		ret.WriteString(numberToWordsHelper(num))
	}
	return ret.String()
}

func numberToWordsHelper(num int) string {
	var ret strings.Builder
	if num >= 100 {
		ret.WriteString(low[num/100])
		ret.WriteString(" Hundred")
	}
	mod := num % 100
	if mod > 0 && ret.Len() > 0 {
		ret.WriteString(" ")
	}
	if mod < 10 {
		ret.WriteString(low[mod])
	} else if mod < 20 {
		ret.WriteString(mid[mod-10])
	} else {
		ret.WriteString(high[mod/10])
		if mod%10 > 0 {
			ret.WriteString(" ")
			ret.WriteString(low[mod%10])
		}
	}

	return ret.String()
}
