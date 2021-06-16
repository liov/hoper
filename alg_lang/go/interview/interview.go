package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func deferCall() {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	panic("触发异常")
}

func goroutine() {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("i: ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("i: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

//组合继承
type People struct{}

func (p *People) ShowA() {
	fmt.Println("showA")
	p.ShowB()
}
func (p *People) ShowB() {
	fmt.Println("showB")
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	fmt.Println("teacher showB")
}

//随机case
func randCase() {
	runtime.GOMAXPROCS(1)
	intChan := make(chan int, 1)
	stringChan := make(chan string, 1)
	intChan <- 1
	stringChan <- "hello"
	select {
	case value := <-intChan:
		fmt.Println(value)
	case value := <-stringChan:
		panic(value)
	}
}

//defer初始化参数
func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func deferInitParam() {
	a := 1
	b := 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	b = 1
}

//并发安全
type UserAges struct {
	ages map[string]int
	sync.Mutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()
	ua.ages[name] = age
}

func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

//数据nil但interface不nil
//interface 类型的变量只有在类型和值均为 nil 时才为 nil
type HuMan interface {
	Show()
}

type Student struct{}

func (stu *Student) Show() {

}

func live() HuMan {
	var stu *Student
	return stu
}

func InterfaceNotNil() {
	if live() == nil {
		fmt.Println("AAAAAAA")
	} else {
		fmt.Println("BBBBBBB")
	}
}

//不初始化的Map是nil
type Param map[string]interface{}

type Show struct {
	Param
}

func NilMap() {
	s := new(Show)
	s.Param["RMB"] = 10000
}

//不确定类型就使用了类型内部字段
type student struct {
	Name string
}

/*func zhoujielun(v interface{}) {
	switch msg := v.(type) {
	case *student, student:
		msg.Name
	}
}*/

//反射无法作用私有字段
type People1 struct {
	name string `json:"name"`
}

func Private() {
	js := `{
		"name":"11"
	}`
	var p People1
	err := json.Unmarshal([]byte(js), &p)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("people: ", p)
}

//%v打印值，打印过程会调用String(),stack overflow
//占位符参数导致对“String”方法的递归调用(%v)
type People2 struct {
	Name string
}

func (p *People2) String() string {
	return fmt.Sprintf("print: %v", p)
}

func overflow() {
	p := &People2{}
	p.String()
}

//提前关闭chan阻塞
func closeChan() {
	ch := make(chan int, 1000)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for {
			a, ok := <-ch
			if !ok {
				fmt.Println("close")
				return
			}
			fmt.Println("a: ", a)
		}
	}()
	close(ch)
	fmt.Println("ok")
	time.Sleep(time.Second * 100)
}

//interface {} is string, not int
type Project struct{}

func (p *Project) deferError() {
	if err := recover(); err != nil {
		fmt.Println("recover: ", err)
	}
}

func (p *Project) exec(msgchan chan interface{}) {
	for msg := range msgchan {
		m := msg.(int)
		fmt.Println("msg: ", m)
	}
}

func (p *Project) run(msgchan chan interface{}) {
	for {
		defer p.deferError()
		go p.exec(msgchan)
		time.Sleep(time.Second * 2)
	}
}

func (p *Project) Main() {
	a := make(chan interface{}, 100)
	go p.run(a)
	go func() {
		for {
			a <- "1"
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(time.Second * 1000000000)
}

func wuyu() {
	p := new(Project)
	p.Main()
}

//当一个被关闭的channel中已经发送的数据都被成功接收后，后续的接收操作将不再阻塞，它们会立即返回一个零值。
// 关闭上面例子中的naturals变量对应的channel并不能终止循环，它依然会收到一个永无休止的零值序列，然后将它们发送给打印者goroutine。
func infiniteLoop() {
	abc := make(chan int, 1000)
	for i := 0; i < 10; i++ {
		abc <- i
	}
	go func() {
		for {
			a := <-abc
			fmt.Println("a: ", a)
		}
	}()
	close(abc)
	fmt.Println("close")
	time.Sleep(time.Second * 100)
}

//在搜索重复时依旧每次都起一个 goroutine 去处理，每个 goroutine 都把它的搜索结果发送到结果 channel 中，channel 中收到的第一条数据会直接返回。
//
//返回完第一条数据后，其他 goroutine 的搜索结果怎么处理？他们自己的协程如何处理？
//
//在 First() 中的结果 channel 是无缓冲的，这意味着只有第一个 goroutine 能返回，由于没有 receiver，其他的 goroutine 会在发送上一直阻塞。如果你大量调用，则可能造成资源泄露。
//
//为避免泄露，你应该确保所有的 goroutine 都能正确退出，有 2 个解决方法：
//
//使用带缓冲的 channel，确保能接收全部 goroutine 的返回结果：
//使用 select 语句，配合能保存一个缓冲值的 channel default 语句：
//default 的缓冲 channel 保证了即使结果 channel 收不到数据，也不会阻塞 goroutine
//使用特殊的废弃（cancellation） channel 来中断剩余 goroutine 的执行：
type query func(string) string

func exec(name string, vs ...query) string {
	ch := make(chan string)
	fn := func(i int) {
		ch <- vs[i](name)
	}
	for i := range vs {
		go fn(i)
	}
	return <-ch
}

func kanbudong() {
	ret := exec("111", func(n string) string {
		return n + "func1"
	}, func(n string) string {
		return n + "func2"
	}, func(n string) string {
		return n + "func3"
	}, func(n string) string {
		return n + "func4"
	})
	fmt.Println(ret)
}

//值接收这种水题也会问吗
type Girl struct {
	Name       string `json:"name"`
	DressColor string `json:"dress_color"`
}

func (g Girl) SetColor(color string) {
	g.DressColor = color
}
func (g Girl) JSON() string {
	data, _ := json.Marshal(&g)
	return string(data)
}
func shui2() {
	g := Girl{Name: "menglu"}
	g.SetColor("white")
	fmt.Println(g.JSON())
}

//切片共享底层数组
func share() {
	str1 := []string{"a", "b", "c"}
	str2 := str1[1:]
	str2[1] = "new"
	fmt.Println(str1)
	str2 = append(str2, "z", "x", "y")
	fmt.Println(str2)
}

type Student3 struct {
	Name string
}
