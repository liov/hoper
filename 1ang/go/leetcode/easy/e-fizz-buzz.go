package easy

import "strconv"

/*
412. Fizz Buzz
写一个程序，输出从 1 到 n 数字的字符串表示。

1. 如果 n 是3的倍数，输出“Fizz”；

2. 如果 n 是5的倍数，输出“Buzz”；

3.如果 n 同时是3和5的倍数，输出 “FizzBuzz”。

示例：

n = 15,

返回:
[

	"1",
	"2",
	"Fizz",
	"4",
	"Buzz",
	"Fizz",
	"7",
	"8",
	"Fizz",
	"Buzz",
	"11",
	"Fizz",
	"13",
	"14",
	"FizzBuzz"

]
*/
func fizzBuzz(n int) []string {
	ret := make([]string, 0, n)
	tmp := 1
	for i := 0; i < n; i++ {
		if tmp == 3 || tmp == 6 || tmp == 9 || tmp == 12 {
			ret = append(ret, "Fizz")
		} else if tmp == 5 || tmp == 10 {
			ret = append(ret, "Buzz")
		} else if tmp == 15 {
			ret = append(ret, "FizzBuzz")
			tmp = 0
		} else {
			ret = append(ret, strconv.Itoa(i+1))
		}
		tmp++
	}
	return ret
}
