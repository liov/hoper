package main

import (
	"fmt"
	"reflect"
)

/*规则一：如果S包含嵌入字段T，则S和*S的方法集都包括具有接收器T的提升方法。*S的方法集还包括具有接收器*T的提升方法。

规则二：如果S包含嵌入字段*T，则S和*S的方法集都包括具有接收器T或*T的提升方法。*/

type People struct {
	Age    int
	gender string
	Name   string
}

type OtherPeople struct {
	People
}

type NewPeople People

func (p *NewPeople) PeName(pname string) {
	fmt.Println("pold name:", p.Name)
	p.Name = pname
	fmt.Println("pnew name:", p.Name)
}

func (p NewPeople) PeInfo() {
	fmt.Println("NewPeople ", p.Name, ": ", p.Age, "岁, 性别:", p.gender)
}

func (p *People) PeName(pname string) {
	fmt.Println("old name:", p.Name)
	p.Name = pname
	fmt.Println("new name:", p.Name)
}

func (p People) PeInfo() {
	fmt.Println("People ", p.Name, ": ", p.Age, "岁, 性别:", p.gender)
}

func methodSet(a interface{}) {
	t := reflect.TypeOf(a)
	fmt.Printf("%T\n", a)
	for i, n := 0, t.NumMethod(); i < n; i++ {
		m := t.Method(i)
		fmt.Println(i, ":", m.Name, m.Type)
	}
}

func main() {
	p := OtherPeople{People{26, "Male", "张三"}}
	p.PeInfo()
	p.PeName("Joke")

	methodSet(p) // T方法提升

	methodSet(&p) // *T和T方法提升

	pp := NewPeople{42, "Male", "李四"}
	pp.PeInfo()
	pp.PeName("Haw")

	methodSet(&pp)
}
