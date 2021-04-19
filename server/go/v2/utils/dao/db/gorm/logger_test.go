package gormi

import (
	"gorm.io/gorm"
	"testing"
)

type Letter int

func (l Letter) String() string {
	switch l {
	case A:
		return "A"
	case B:
		return "B"
	default:
		return "C"
	}
}

const (
	A Letter = 1 + iota
	B
	C
)

type Foo struct {
	Id     int
	Letter Letter
}

func TestLogger(t *testing.T) {
	db := gorm.DB{}
	var foos []Foo
	db.Where("letter=?", A).Find(&foos)
}
