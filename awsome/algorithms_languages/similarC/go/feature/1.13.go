package main

import (
	"errors"
	"fmt"
)

//Go设计者决定在Go 1.13版本中增加Go对数字字面量的表达能力，在这方面对Go语言做了如下补充：
//
//增加二进制数字字面量，以0b或0B开头
//
//在保留以”0″开头的八进制数字字面量形式的同时，增加以”0o”或”0O”开头的八进制数字字面量形式
//
//增加十六进制形式的浮点数字面量，以0x或0X开头的、形式如0×123.86p+2的浮点数
//
//为提升可读性，在数字字面量中增加数字分隔符”_”，分隔符可以用来分隔数字(起到分组提高可读性作用，比如每3个数字一组)，也可以用来分隔前缀与第一个数字。
func number() {
	var a = 0b101
	var b = 0o123
	var c = 0x123
	var d = 0x0.1p7
	var e = 6 + 5i
	var f = 1_0000_0000
	fmt.Println(a, b, c, d, e, f)
	var g int = 5
	fmt.Println(5 << g)
}

//Go 1.13中关于语言规范方面的另一个变动点是取消了移位操作(>>的<<)的右操作数仅能是无符号数的限制，以前必须的强制到uint的转换现在不必要了：
func displacement() {
	var i int = 5
	fmt.Println(2 << uint(i)) // before go 1.13
	fmt.Println(2 << i)
}

//面向私有模块的GOPRIVATE
//有了GOPROXY后，公共module的数据获取变得十分easy。但是如果依赖的是企业内部module或托管站点上的private库，通过GOPROXY（默认值）获取显然会得到一个失败的结果，除非你搭建了自己的公私均可的goproxy server并将其设置到GOPROXY中。
//
//Go 1.13提供了GOPRIVATE变量，用于指示哪些仓库下的module是private，不需要通过GOPROXY下载，也不需要通过GOSUMDB去验证其校验和。不过要注意的是GONOPROXY和GONOSUMDB可以override GOPRIVATE中的设置，因此设置时要谨慎，比如下面的例子：

//通过标准库增加了errors.Is和As函数来解决error value比较问题
//
//增加errors.Unwrap来解决error unwrap问题。
func err() error {
	return errors.Unwrap(fmt.Errorf("%w", errors.New("UnWarp")))
}
