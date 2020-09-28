package main

import "fmt"

//命令模式是将一个请求封装为一个对象,从而使我们可用不同的请求对客户进行
//参数化;对请求排队或者记录请求日志,以及支持可撤销的操作.命令模式是一种对象行
//为型模式,其别名为动作模式或事务模式.
/*客户端通过调用者发送命令,命令调用接收者执行相应操作.
其实命令模式也很简单,不过不知道大家发现没有,在上述描述中调用者和接收者并不知道对方的
存在,也就是说他们之间是解耦合的.
还是用遥控器的例子来解释一下吧,遥控器对应上面的角色就是调用者,电视就是接收者,命令呢?
对应的就是遥控器上的按键,最后客户端就是我们人啦,当我们想打开电视的时候,就会通过遥控器(调用者)
的电源按钮(命令)来打开电视(接收者),在这个过程中遥控器是不知道电视的,但是电源按钮是知道他要控制
谁的什么操作.*/

// receiver
type TV struct{}

func (p *TV) Open() {
	fmt.Println("play...")
}
func (p *TV) Close() {
	fmt.Println("stop...")
}

/*这台电视弱爆了,只有打开和关闭两个功能,对应的就是上面代码中的 Open 和 CloseDao 两个方法,虽
然 很 简 单 , 不 过 我 们 确 实 造 出 了 一 台 电 视 ( 估 计 还 是 彩 色 的 ) .
下面,我们再来实现命令,也就是遥控器上的按键,因为电视只有打开和关闭功能,所以按键我们也只提供两个,
多了用不上.*/
// command
type Command interface {
	Press()
}

type OpenCommand struct {
	tv TV
}

func (p OpenCommand) Press() {
	p.tv.Open()
}

type CloseCommand struct {
	tv TV
}

func (p CloseCommand) Press() {
	p.tv.Close()
}

/*首先我们定义了一个命令接口,只有一个方法就是 Press ,当我们按下按键时会去调用这个方法,然
后我们果然只实现了两个按键,分别是 OpenCommand 和 CloseCommand ,这两个实现中都保存着一台电视的句
柄,并且在 Press 方法中根据功能去调用了这个 tv 的相应方法来完成正确的操作.
还有什么我们没有实现? 调用者,也就是遥控器了,来看看遥控器怎么实现吧.*/
// sender
type Invoker struct {
	cmd Command
}

func (p *Invoker) SetCommand(cmd Command) {
	p.cmd = cmd
}
func (p Invoker) Do() {
	p.cmd.Press()
}
