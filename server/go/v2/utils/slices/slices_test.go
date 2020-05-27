package slices

import (
	"fmt"
	"log"
	"testing"
)

func TestContains(t *testing.T) {
	val1 := []string{"a", "b", "c"}
	val2 := "a"
	val3 := "d"
	fmt.Println(StringContains(val1, val2))
	fmt.Println(StringContains(val1, val3))
}

type Foo struct {
	ID  uint64
	Str string
}

func (f *Foo) IsEqual(v interface{}) bool {
	if f1, ok := v.(*Foo); ok {
		if f1.ID == f.ID {
			return true
		}
	}
	return false
}

var _ Equal = &Foo{}

func TestIsEqu(t *testing.T) {
	s1 := []Equal{
		&Foo{1, "1"},
		&Foo{2, "2"},
		&Foo{3, "3"},
	}
	s2 := []Equal{
		&Foo{4, "1"},
		&Foo{5, "1"},
		&Foo{6, "1"},
	}
	log.Println(IsCoincide(s1, s2))
}
