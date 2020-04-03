package aop

import (
	"log"

	"bou.ke/monkey"
)

func init() {
	monkey.Patch(NewDao, func() Dao {
		return &UserProxy{&User{}}
	})
}

type Dao interface {
	Save()
	Update()
}

type User struct {
}

//go:noinline
func NewDao() Dao {
	return &User{}
}

func (u *User) Save() {

}
func (u *User) Update() {

}

type UserProxy struct {
	user *User
}

func (u *UserProxy) Save() {
	log.Println("save")
	u.user.Save()
}
func (u *UserProxy) Update() {
	log.Println("update")
	u.user.Update()
}

type Proxy struct {
	p interface{}
}
