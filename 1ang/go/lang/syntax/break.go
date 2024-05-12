package main

import (
	"fmt"
	"time"
)

/*Go语言select、switch中使用break有用吗？
问题-select中使用break有用吗？
*/

func test3() {

SELECT:
	for {
		select {
		case <-time.After(time.Second):
			fmt.Println("一秒后退出")
			//break 跳出select
			break SELECT //带标签的break，实际上跳出到select外层的for语句块
		case <-time.After(time.Second * 10):
			fmt.Println("十秒后退出")
			break
		}
	}

	fmt.Println("select 语句结束")
}

/**
output:
一秒后退出
select 语句结束
*/

/*以上例子可以看出：

带标签的break，可以跳出多层select/ switch作用域。让break更加灵活，写法更加简单灵活，不需要使用控制变量一层一层跳出循环，没有带break的只能跳出当前语句块。



/*通过以上示例可以得出

没有fallthough的switch语句break也是没啥用的

示例2 有fallthough的例子-一般情况
*/

func switchWithoutFallthrough(i int) {
	switch i {
	case 10:
		fmt.Println("等于10")
		break
	case 8:
		fmt.Println("等于8")
		break
	case 5:
		fmt.Println("等于5")
		break
	default:
		fmt.Println("不关心")
		break
	}
}

func switchWithFallthrough(i int) {
	switch {
	case i < 10:
		fmt.Println("等于10")
		break
	case i < 8:
		fmt.Println("等于8")
		fallthrough
	case i < 5:
		fmt.Println("等于5")
		break
	default:
		fmt.Println("不关心")
		break
	}
}
func main() {
	test3()

	switchWithoutFallthrough(8)
	switchWithFallthrough(8)

}

/**
output:
等于8
等于5
等于8
*/

/*以上例子可以看出

增加fallthrough以后会执行fallthrough后面的case语句，通过break可以跳过fallthrough的顺序执行。这种情况下break是有用的

通过select和break的分析，我们也可以推断出 带标签的break，可以跳出多层switch作用域。让break更加灵活，写法更加简单灵活，不需要使用控制变量一层一层跳出循环，没有带break的只能跳出当前语句块。（此处例子可以自行实现）

总结
1. 单独在select中使用break和不使用break没有啥区别。
2. 单独在表达式switch语句，并且没有fallthough，使用break和不使用break没有啥区别。
3. 单独在表达式switch语句，并且有fallthough，使用break能够终止fallthough后面的case语句的执行。
4. 带标签的break，可以跳出多层select/ switch作用域。让break更加灵活，写法更加简单灵活，不需要使用控制变量一层一层跳出循环，没有带break的只能跳出当前语句块。
*/
