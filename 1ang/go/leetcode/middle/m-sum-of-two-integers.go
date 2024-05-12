package middle

/*
371. 两整数之和
给你两个整数 a 和 b ，不使用 运算符 + 和 - ​​​​​​​，计算并返回两整数之和。



示例 1：

输入：a = 1, b = 2
输出：3
示例 2：

输入：a = 2, b = 3
输出：5


提示：

-1000 <= a, b <= 1000
https://leetcode-cn.com/problems/sum-of-two-integers/
*/

func getSum(a int, b int) int {
	sum := a ^ b
	carry := a & b << 1
	for carry != 0 {
		tmp := sum
		sum = sum ^ carry
		carry = tmp & carry << 1
	}
	return sum
}
