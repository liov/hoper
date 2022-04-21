package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func main() {
	u := User{Id: 1}
	u1 := User{Id: 2}
	compare([]*User{&u1, &u})
	//compare1([]*User{&u1, &u}) //Cannot use '[]*User{&u1, &u}' (type []*User) as the type []Compared
}

type Queue interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func name[T constraints.Integer](x T) {
	fmt.Println(x)
}

type Compared interface {
	Less(t Compared) bool
}

type User struct {
	Id int
}

func (u *User) Less(u1 Compared) bool {
	return u.Id < u1.(*User).Id
}

func compare[T Compared](t []T) {
	fmt.Println(t[0].Less(t[1]))
}

//错误写法
func compare1(t []Compared) {
	fmt.Println(t[0].Less(t[1]))
}
