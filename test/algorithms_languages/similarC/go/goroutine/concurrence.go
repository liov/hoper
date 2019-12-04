package main

//并发控制
/*1. 使用最基本通过channel通知实现并发控制
无缓冲通道
无缓冲的通道指的是通道的大小为0，也就是说，这种类型的通道在接收前没有能力保存任何值，它要求发送 goroutine 和接收 goroutine 同时准备好，才可以完成发送和接收操作。
从上面无缓冲的通道定义来看，发送 goroutine 和接收 gouroutine 必须是同步的，同时准备后，如果没有同时准备好的话，先执行的操作就会阻塞等待，直到另一个相对应的操作准备好为止。这种无缓冲的通道我们也称之为同步通道。
正式通过无缓冲通道来实现多 goroutine 并发控制
2.在 sync 包中，提供了 WaitGroup ，它会等待它收集的所有 goroutine 任务全部完成。在WaitGroup里主要有三个方法

Add, 可以添加或减少 goroutine的数量
Done, 相当于Add(-1)
Wait, 执行后会堵塞主线程，直到WaitGroup 里的值减至0

在主 goroutine 中 Add(delta int) 索要等待goroutine 的数量。在每一个 goroutine 完成后 Done() 表示这一个goroutine 已经完成，当所有的 goroutine 都完成后，在主 goroutine 中 WaitGroup 返回返回。

3. 在Go 1.7 以后引进的强大的Context上下文，实现并发控制

在一些简单场景下使用 channel 和 WaitGroup 已经足够了，但是当面临一些复杂多变的网络并发场景下 channel 和 WaitGroup 显得有些力不从心了。比如一个网络请求 Request，每个 Request 都需要开启一个 goroutine 做一些事情，这些 goroutine 又可能会开启其他的 goroutine，比如数据库和RPC服务。所以我们需要一种可以跟踪 goroutine 的方案，才可以达到控制他们的目的，这就是Go语言为我们提供的 Context，称之为上下文非常贴切，它就是goroutine 的上下文。它是包括一个程序的运行环境、现场和快照等。每个程序要运行时，都需要知道当前程序的运行状态，通常Go 将这些封装在一个 Context 里，再将它传给要执行的 goroutine 。
context 包主要是用来处理多个 goroutine 之间共享数据，及多个 goroutine 的管理。

*/

/*并发(concurrency)：两个或两个以上的任务在一段时间内被执行。我们不必care这些任务在某一个时间点是否是同时执行，可能同时执行，也可能不是，我们只关心在一段时间内，哪怕是很短的时间（一秒或者两秒）是否执行解决了两个或两个以上任务。

并行(parallellism)：两个或两个以上的任务在同一时刻被同时执行。

并发说的是逻辑上的概念，而并行，强调的是物理运行状态。并发“包含”并行。

（详情请见：Rob Pike 的PPT）

Go的CSP并发模型
Go实现了两种并发形式。第一种是大家普遍认知的：多线程共享内存。其实就是Java或者C++等语言中的多线程开发。另外一种是Go语言特有的，也是Go语言推荐的：CSP（communicating sequential processes）并发模型。

CSP并发模型是在1970年左右提出的概念，属于比较新的概念，不同于传统的多线程通过共享内存来通信，CSP讲究的是“以通信的方式来共享内存”。

请记住下面这句话：
Do not communicate by sharing memory; instead, share memory by communicating.
“不要以共享内存的方式来通信，相反，要通过通信来共享内存。”

普通的线程并发模型，就是像Java、C++、或者Python，他们线程间通信都是通过共享内存的方式来进行的。非常典型的方式就是，在访问共享数据（例如数组、Map、或者某个结构体或对象）的时候，通过锁来访问，因此，在很多时候，衍生出一种方便操作的数据结构，叫做“线程安全的数据结构”。例如Java提供的包”java.util.concurrent”中的数据结构。Go中也实现了传统的线程并发模型。

Go的CSP并发模型，是通过goroutine和channel来实现的。

goroutine 是Go语言中并发的执行单位。有点抽象，其实就是和传统概念上的”线程“类似，可以理解为”线程“。
channel是Go语言中各个并发结构体(goroutine)之前的通信机制。 通俗的讲，就是各个goroutine之间通信的”管道“，有点类似于Linux中的管道。*/

/*Go线程实现模型MPG
M指的是Machine，一个M直接关联了一个内核线程。
P指的是”processor”，代表了M所需的上下文环境，也是处理用户级代码逻辑的处理器。
G指的是Goroutine，其实本质上也是一种轻量级的线程。
一个M会对应一个内核线程，一个M也会连接一个上下文P，一个上下文P相当于一个“处理器”，一个上下文连接一个或者多个Goroutine。P(Processor)的数量是在启动时被设置为环境变量GOMAXPROCS的值，或者通过运行时调用函数runtime.GOMAXPROCS()进行设置。Processor数量固定意味着任意时刻只有固定数量的线程在运行go代码。Goroutine中就是我们要执行并发的代码。图中P正在执行的Goroutine为蓝色的；处于待执行状态的Goroutine为灰色的，灰色的Goroutine形成了一个队列runqueues
我们能不能直接除去上下文，让Goroutine的runqueues挂到M上呢？答案是不行，需要上下文的目的，是让我们可以直接放开其他线程，当遇到内核线程阻塞的时候。

一个很简单的例子就是系统调用sysall，一个线程肯定不能同时执行代码和系统调用被阻塞，这个时候，此线程M需要放弃当前的上下文环境P，以便可以让其他的Goroutine被调度执行。
均衡的分配工作
按照以上的说法，上下文P会定期的检查全局的goroutine 队列中的goroutine，以便自己在消费掉自身Goroutine队列的时候有事可做。假如全局goroutine队列中的goroutine也没了呢？就从其他运行的中的P的runqueue里偷。

每个P中的Goroutine不同导致他们运行的效率和时间也不同，在一个有很多P和M的环境中，不能让一个P跑完自身的Goroutine就没事可做了，因为或许其他的P有很长的goroutine队列要跑，得需要均衡。

stack
OS线程初始栈为2MB。Go语言中，每个goroutine采用动态扩容方式，初始2KB，按需增长，最大1G。此外GC会收缩栈空间。
BTW，增长扩容都是有代价的，需要copy数据到新的stack，所以初始2KB可能有些性能问题。
更多关于stack的内容，可以参见大佬的文章。聊一聊goroutine stack
管理
用户线程的调度以及生命周期管理都是用户层面，Go语言自己实现的，不借助OS系统调用，减少系统资源消耗。
G-M-P
Go语言采用两级线程模型，即用户线程与内核线程KSE（kernel scheduling entity）是M:N的。最终goroutine还是会交给OS线程执行，但是需要一个中介，提供上下文。这就是G-M-P模型

G: goroutine, 类似进程控制块，保存栈，状态，id，函数等信息。G只有绑定到P才可以被调度。
M: machine, OS线程，绑定有效的P之后，进入调度。
P: 逻辑处理器，保存各种队列G。对于G而言，P就是cpu core。对于M而言，P就是上下文。P的数量由GOMAXPROCS设置，最大256。
sched: 调度程序，保存GRQ，midle M空闲队列，pidle P空闲队列以及lock等信息






G-M-P模型

队列
Go调度器有两个不同的运行队列：

GRQ，全局运行队列，尚未分配给P的G
LRQ，本地运行队列，每个P都有一个LRQ，用于管理分配给P执行的G

状态
go1.10\src\runtime\runtime2.go

_Gidle: 分配了G，但是没有初始化
_Grunnable: 在run queue运行队列中，LRQ或者GRQ
_Grunning: 正在运行指令，有自己的stack。不在runq运行队列中，分配给M和P
_Gsyscall: 正在执行syscall，而非用户指令，不在runq，分给M，P给找idle的M
_Gwaiting: block。不在RQ，但是可能会在channel的wait queue等待队列
_Gdead: unused。在P的gfree list中，不在runq。idle闲置状态
_Gcopystack: stack扩容或者gc收缩

上下文切换
Go调度器根据事件进行上下文切换。

go关键字，创建goroutine
gc垃圾回收，gc也是goroutine，所以需要时间片
system call系统调用，block当前G
sync同步，block当前G

调度
调度的目的就是防止M堵塞，空闲，系统进程切换。
详见 Golang - 调度剖析【第二部分】
异步调用
Linux可以通过epoll实现网络调用，统称网络轮询器N（Net Poller）。

G1在M上运行，P的LRQ有其他3个G，N空闲；
G1进行网络IO，因此被移动到N，M继续从LRQ取其他的G执行。比如G2就被上下文切换到M上；
G1结束网络请求，收到响应，G1被移回LRQ，等待切换到M执行。

同步调用
文件IO操作

G1在M1上运行，P的LRQ有其他3个G；
G1进行同步调用，堵塞M；
调度器将M1与P分离，此时M1下只有G1，没有P。
将P与空闲M2绑定，并从LRQ选择G2切换
G1结束堵塞操作，移回LRQ。M1空闲备用。

任务窃取
上面都是防止M堵塞，任务窃取是防止M空闲

两个P，P1，P2
如果P1的G都执行完了，LRQ空，P1就开始任务窃取。
第一种情况，P2 LRQ还有G，则P1从P2窃取了LRQ中一半的G
第二种情况，P2也没有LRQ，P1从GRQ窃取。

g0
每个M都有一个特殊的G，g0。用于执行调度，gc，栈管理等任务，所以g0的栈称为调度栈。g0的栈不会自动增长，不会被gc，来自os线程的栈。

*/
