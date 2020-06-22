简单描述下优化的原理，go 1.14新加入了开放编码（Open-coded）defer类型，编译器在ssa过程中会把被延迟的方法直接插入到函数的尾部，避免了运行时的deferproc及deferprocStack操作。
免除了在没有运行时判断下的deferreturn调用。
如有运行时判断的逻辑，则deferreturn也进一步优化，开放编码下的deferreturn不会进行jmpdefer的尾递归调用，而直接在一个循环里遍历执行。

共有三种defer模式类型，编译后一个函数里只会一种defer模式。

** 第一种，堆上分配 (deferProc)，基本是依赖运行时来分配 “_defer” 对象并加入延迟参数。在函数的尾部插入deferreturn方法来消费defer链条。

** 第二种，栈上分配 (deferprocStack)，基本跟堆上差不多，只是分配方式改为在栈上分配，压入的函数调用栈存有_defer记录，另外编译器在ssa过程中会预留defer空间。

** 第三种，open coded开放编码模式。open coded就是go 1.14新增的模式。

默认open-coded最多支持8个defer，超过则取消。

在构建ssa时如发现gcflags有N禁止优化的参数 或者 return数量 * defer数量超过了 15不适用open-coded模式。

逃逸分析会判断循序的层数，如果有轮询，那么强制使用栈分配模式。

ssa的过程？延迟比特又是什么？

ssa的构建过程是相当的麻烦，源码也很难理解。简单说golang编译器可在中间代码ssa的过程中优化defer，open-coded模式下把被延迟的方法和deferreturn直接插入到函数尾部。

多数defer看直接被编译器分析优化，但如果一个 defer 发生在一个条件语句中，而这个条件必须等到运行时才能确定。
go1.14在open-coded使用延迟比特 (defer bit) 来判断条件分支是否该执行。一个字节8个比特，在open coded里最多8个defer，包括了if判断里的defer。
只要有defer关键字就在相应位置设置bit，而在if判断里需要运行时来设置bit，编译器无法控制的。

需要运行时确认defer的例子