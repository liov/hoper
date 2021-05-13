package slices

import (
	"github.com/liov/hoper/v2/utils/def"
	"log"
	"testing"
)

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
func (f *Foo) CmpKey() uint64 {
	return f.ID
}

var _ Equal = &Foo{}

func TestIsEqu(t *testing.T) {
	s1 := []def.CmpKey{
		&Foo{1, "1"},
		&Foo{2, "2"},
		&Foo{3, "3"},
	}
	s2 := []def.CmpKey{
		&Foo{4, "1"},
		&Foo{5, "1"},
		&Foo{6, "1"},
	}
	log.Println(IsCoincide(s1, s2))
}

func TestDiff(t *testing.T) {
	a := []uint64{1, 2, 3, 4}
	b := []uint64{2, 3, 4, 5}
	t.Log(DiffUint64(a, b))
}
