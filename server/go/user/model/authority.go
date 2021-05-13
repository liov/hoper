package model

import (
	"fmt"
	"strings"
)

//权限
type Authority uint64

const (
	Create Authority = 1 << iota
	EditUser
	DeleteUser
)

func (a Authority) String() string {
	switch a {
	case DeleteUser:
		return "删除用户"
	}
	return ""
}

var authorities = []Authority{Create, EditUser, DeleteUser}

func (a Authority) Parsed() []Authority {
	var retAuthorities []Authority
	for i := range authorities {
		if authorities[i]&a == a {
			retAuthorities = append(retAuthorities, authorities[i])
		}
	}
	return retAuthorities
}

type Authorities []Authority

func (f Authorities) Sql() string {
	var sqls []string
	for i := range f {
		sqls = append(sqls, fmt.Sprintf("authority&%d = %d", f[i], f[i]))
	}
	return strings.Join(sqls, " || ")
}

func (f Authorities) Value() int {
	var ret int
	for i := range f {
		ret |= int(f[i])
	}
	return ret
}

type Hobby uint32

const (
	Douyin Hobby = iota
	Dance
	Sing
)

var hobbies = []Hobby{Douyin, Dance, Sing}

func (f Hobby) Parsed() []Hobby {
	var retHobbies []Hobby
	for i := range hobbies {
		if 1<<hobbies[i]&f == 1<<hobbies[i] {
			retHobbies = append(retHobbies, hobbies[i])
		}
	}
	return retHobbies
}

type Hobbys []Hobby

func (f Hobbys) Sql() string {
	var sqls []string
	for i := range f {
		sqls = append(sqls, fmt.Sprintf("hobby&%d = %d", 1<<f[i], 1<<f[i]))
	}
	return strings.Join(sqls, " || ")
}

func (f Hobbys) Value() int {
	var ret int
	for i := range f {
		ret |= 1 << f[i]
	}
	return ret
}
