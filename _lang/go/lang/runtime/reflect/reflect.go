package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type A int

// 普通函数
func (*A) M1(a, b int) int {
	return a + b
}

func m2(a, b int) int {
	return a + b
}

type MyType struct {
	i    int
	name string
}

func (mt *MyType) SetI(i int) {
	mt.i = i
}

func (mt *MyType) SetName(name string) {
	mt.name = name
}

func (mt *MyType) String() string {
	return fmt.Sprintf("%p", mt) + "--name:" + mt.name + " i:" + strconv.Itoa(mt.i)
}

func main() {
	var a A
	// 取变量a的反射类型对象
	typeOfA := reflect.TypeOf(a)
	// 根据反射类型对象创建类型实例
	aIns := reflect.New(typeOfA)
	// 输出Value的类型和种类
	fmt.Println(aIns.Type(), aIns.Kind())

	var b = A(1)
	//方法要公开，要传入指针且不Elem()，用TypeOf会报参数不足，要传入对象，just so so
	m1, _ := reflect.TypeOf(&b).MethodByName("M1")
	params1 := make([]reflect.Value, 3)
	params1[0] = reflect.ValueOf(&b)
	params1[1] = reflect.ValueOf(18)
	params1[2] = reflect.ValueOf(12)
	mv := m1.Func.Call(params1)
	// 获取第一个返回值, 取整数值
	fmt.Println("m1:", mv[0].Int())

	// 将函数包装为反射值对象
	funcValue := reflect.ValueOf(m2)
	// 构造函数参数, 传入两个整型值
	paramList := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
	// 反射调用函数
	retList := funcValue.Call(paramList)
	// 获取第一个返回值, 取整数值
	fmt.Println("m2:", retList[0].Int())

	//真的很奇怪，对指针取指针
	myType := &MyType{22, "golang"}
	params := make([]reflect.Value, 1)
	mtV := reflect.ValueOf(myType)
	m2 := mtV.MethodByName("String")
	fmt.Println("Before:", m2.Call(nil)[0])
	params[0] = reflect.ValueOf(18)
	mtV.MethodByName("SetI").Call(params)
	params[0] = reflect.ValueOf("reflection hoper")
	mtV.MethodByName("SetName").Call(params)
	fmt.Println("After:", m2.Call(nil)[0])
}
