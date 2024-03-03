package easy

/*
367. 有效的完全平方数
给定一个 正整数 num ，编写一个函数，如果 num 是一个完全平方数，则返回 true ，否则返回 false 。

进阶：不要 使用任何内置的库函数，如  sqrt 。

示例 1：

输入：num = 16
输出：true
示例 2：

输入：num = 14
输出：false

提示：

1 <= num <= 2^31 - 1

https://leetcode-cn.com/problems/valid-perfect-square/
*/
func isPerfectSquare(num int) bool {
	_, ok := square[num]
	return ok
}

var square = make(map[int]struct{})

func init() {

	for i := 0; i < 1<<16; i++ {
		sqr := i * i
		if sqr < 1<<31-1 {
			square[sqr] = struct{}{}
		} else {
			break
		}
	}
}
