package main

import "fmt"

type Car struct {
	name string
}

type Cars []*Car

func (cs Cars) Process(f func(car *Car)) {
	for _, c := range cs {
		f(c)
	}
}

var G int = 7

func main() {

	fn := func() {
		fmt.Println("hello")
	}
	fn()

	fmt.Println("匿名函数加法求和：", func(x, y int) int { return x + y }(3, 4))

	func() {
		sum := 0
		for i := 1; i <= 1e6; i++ {
			sum += i
		}
		fmt.Println("匿名函数加法循环求和：", sum)
	}()

	// 影响全局变量G，代码块状态持续
	y := func() int {
		fmt.Printf("G: %d, G的地址:%p\n", G, &G)
		G += 1
		return G
	}
	fmt.Println(y(), y)
	fmt.Println(y(), y)
	fmt.Println(y(), y) //y的地址

	// 影响全局变量G，注意z的匿名函数是直接执行，所以结果不变
	z := func() int {
		G += 1
		return G
	}()
	fmt.Println(z, &z)
	fmt.Println(z, &z)
	fmt.Println(z, &z)

	// 影响外层（自由）变量i，代码块状态持续
	var f = N()
	fmt.Println(f(1), &f)
	fmt.Println(f(1), &f)
	fmt.Println(f(1), &f)

	var f1 = N()
	fmt.Println(f1(1), &f1)

}

func N() func(int) int {
	var i int
	return func(d int) int {
		fmt.Printf("i: %d, i的地址:%p\n", i, &i)
		i += d
		return i
	}
}
