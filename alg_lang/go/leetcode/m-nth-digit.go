package leetcode

import "math"

/*
400. 第 N 位数字
给你一个整数 n ，请你在无限的整数序列 [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ...] 中找出并返回第 n 位上的数字。



示例 1：

输入：n = 3
输出：3
示例 2：

输入：n = 11
输出：0
解释：第 11 位数字在序列 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ... 里是 0 ，它是 10 的一部分。


提示：

1 <= n <= 2^31 - 1
第 n 位上的数字是按计数单位（digit）从前往后数的第 n 个数，参见 示例 2 。

https://leetcode-cn.com/problems/nth-digit/
*/

func findNthDigit(n int) int {
	i := 1
	lastNum := 0
	for {
		nowNum := int(math.Pow10(i))
		tmp := n
		n -= i * (nowNum - 1 - lastNum)
		if n <= 0 {
			num := lastNum + tmp/i
			mod := tmp % i
			if mod > 0 {
				num++
				num = num / int(math.Pow10(i-mod))
			}
			return num % 10
		}
		lastNum = nowNum - 1
		i++
	}
}
