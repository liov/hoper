package pprof

import (
	"net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"
)

//install graphviz
//go tool pprof http://localhost:8080/debug/pprof/profile?seconds=60

//Go性能优化
//Go语言项目中的性能优化主要有以下几个方面：
//
//CPU profile：报告程序的 CPU 使用情况，按照一定频率去采集应用程序在 CPU 和寄存器上面的数据
//Memory Profile（Heap Profile）：报告程序的内存使用情况
//Block Profiling：报告 goroutines 不在运行状态的情况，可以用来分析和查找死锁等性能瓶颈
//Goroutine Profiling：报告 goroutines 的使用情况，有哪些 goroutine，它们的调用关系是怎样的
//采集性能数据
//Go语言内置了获取程序的运行数据的工具，包括以下两个标准库：
//
//runtime/pprof：采集工具型应用运行数据进行分析
//net/http/pprof：采集服务型应用运行时数据进行分析
//pprof开启后，每隔一段时间（10ms）就会收集下当前的堆栈信息，获取格格函数占用的CPU以及内存资源；最后通过对这些采样数据进行分析，形成一个性能分析报告。
//
//注意，我们只应该在性能测试的时候才在代码中引入pprof。
//
//工具型应用
//如果你的应用程序是运行一段时间就结束退出类型。那么最好的办法是在应用退出的时候把 profiling 的报告保存到文件中，进行分析。对于这种情况，可以使用runtime/pprof库。 首先在代码中导入runtime/pprof工具：
//
//import "runtime/pprof"
//CPU性能分析
//开启CPU性能分析：
//
//pprof.StartCPUProfile(w io.Writer)
//停止CPU性能分析：
//
//pprof.StopCPUProfile()
//应用执行结束后，就会生成一个文件，保存了我们的 CPU profiling 数据。得到采样数据之后，使用go tool pprof工具进行CPU性能分析。
//
//内存性能优化
//记录程序的堆栈信息
//
//pprof.WriteHeapProfile(w io.Writer)
//得到采样数据之后，使用go tool pprof工具进行内存性能分析。
//
//go tool pprof默认是使用-inuse_space进行统计，还可以使用-inuse-objects查看分配对象的数量。
func main() {
	http.ListenAndServe(":8080", nil)
}
