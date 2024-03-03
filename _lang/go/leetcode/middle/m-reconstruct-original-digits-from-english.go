package middle

import (
	"strings"
)

/*
423. 从英文中重建数字
给你一个字符串 s ，其中包含字母顺序打乱的用英文单词表示的若干数字（0-9）。按 升序 返回原始的数字。



示例 1：

输入：s = "owoztneoer"
输出："012"
示例 2：

输入：s = "fviefuro"
输出："45"


提示：

1 <= s.length <= 105
s[i] 为 ["e","g","f","i","h","o","n","s","r","u","t","w","v","x","z"] 这些字符之一
s 保证是一个符合题目要求的字符串

https://leetcode-cn.com/problems/reconstruct-original-digits-from-english/
*/

/*
可以发现z，w,u,x,g 都只在一个数字中，即 0,2,4,6,8 中出现。因此我们可以使用一个哈希表统计每个字母出现的次数，那么 z,w,u,x,g 出现的次数，即分别为 0,2,4,6,8 出现的次数。

随后我们可以注意那些只在两个数字中出现的字符：

h 只在 3,8 中出现。由于我们已经知道了 8 出现的次数，因此可以计算出 3 出现的次数。

f 只在 4,5 中出现。由于我们已经知道了 4 出现的次数，因此可以计算出 5 出现的次数。

s 只在 6,7 中出现。由于我们已经知道了 6 出现的次数，因此可以计算出 7 出现的次数。

作者：LeetCode-Solution
链接：https://leetcode-cn.com/problems/reconstruct-original-digits-from-english/solution/cong-ying-wen-zhong-zhong-jian-shu-zi-by-9g1r/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func originalDigits(s string) string {
	statistics := make([]int, 26)
	for i := range s {
		statistics[s[i]-'a']++
	}
	count := make([]int, 10)
	var ans strings.Builder
	count[0] = statistics['z'-'a']
	count[2] = statistics['w'-'a']
	count[4] = statistics['u'-'a']
	count[6] = statistics['x'-'a']
	count[8] = statistics['g'-'a']
	count[3] = statistics['h'-'a'] - count[8]
	count[5] = statistics['f'-'a'] - count[4]
	count[7] = statistics['s'-'a'] - count[6]
	count[1] = statistics['o'-'a'] - count[0] - count[2] - count[4]
	count[9] = statistics['i'-'a'] - count[5] - count[6] - count[8]

	for i, c := range count {
		for j := 0; j < c; j++ {
			ans.WriteByte(numberBytes[i])
		}
	}
	return ans.String()
}
