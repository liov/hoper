package main

import (
	"fmt"
	"reflect"
)

/**
 * @author     ：lbyi
 * @date       ：Created in 2019/3/26
 * @description：
 */
type User struct {
	Name string
	User *User
}

func main() {
	user1 := User{Name: "a"}
	user2 := User{"b", &user1}
	fmt.Print(user1, user2.User)
}

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human  //匿名字段
	school string
	loan   float32
}

type Employee struct {
	Human   //匿名字段
	company string
	money   float32
}

//Human实现SayHi方法
func (h *Human) SayHi() {
	fmt.Printf("嗨,我是%s你可以给我打电话%s\n", h.name, h.phone)
}

//Human实现Sing方法
func (h *Human) Sing(lyrics string) {
	fmt.Println("啦啦啦", lyrics)
}

//Employee重载Human的SayHi方法
func (e *Employee) SayHi() {
	fmt.Printf("嗨, 我是%s,我在%s工作.叫我%s\n", e.name,
		e.company, e.phone)
}

// Interface Men被Human,Student和Employee实现
// 因为这三个类型都实现了这两个方法
type Men interface {
	SayHi()
	Sing(lyrics string)
}

func Ts() {
	mike := Student{Human{"Mike", 25, "222-222-XXX"}, "清华", 0.00}
	paul := Student{Human{"Paul", 26, "111-222-XXX"}, "哈佛", 100}
	sam := Employee{Human{"Sam", 36, "444-222-XXX"}, "谷歌.", 1000}
	tom := Employee{Human{"Tom", 37, "222-444-XXX"}, "汇桔.", 5000}

	//定义Men类型的变量i
	var i Men

	//i能存储Student
	i = &mike
	fmt.Println("这是麦克, 一个学生:")
	i.SayHi()
	i.Sing("11月下雨")

	//i也能存储Employee
	i = &tom
	fmt.Println("这是汤姆，一个职员:")
	i.SayHi()
	i.Sing("布朗是野的")

	//定义了slice Men
	fmt.Println("让我们用一个人切片看会发生什么")
	x := make([]Men, 3)
	//这三个都是不同类型的元素，但是他们实现了interface同一个接口
	x[0], x[1], x[2] = &paul, &sam, &mike

	for _, value := range x {
		value.SayHi()
	}
	t := reflect.TypeOf(&mike)
	t2 := reflect.TypeOf(mike) //得到类型的元数据,通过t我们能获取类型定义里面的所有元素
	//t3:= reflect.TypeOf(mike).String()
	v := reflect.ValueOf(&mike) //得到实际的值，通过v我们获取存储在里面的值，还可以去改变值
	v2 := reflect.ValueOf(mike)
	tag := t.Elem()  //获取定义在struct里面的标签
	name := v.Elem() //获取存储在第一个字段里面的值

	for i := 0; i < t2.NumField(); i++ {
		fmt.Println(t2.Field(i).Name)
	}
	v2.Field(0).Interface()
	fmt.Println(" type:", tag)
	fmt.Println("name:", name)
	fmt.Println("name:", v.Elem().Field(0).Interface())

}
