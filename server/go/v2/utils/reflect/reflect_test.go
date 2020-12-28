package reflecti

import (
	"testing"

	"github.com/liov/hoper/go/v2/utils/log"
)

type Foo struct {
	A int
	B string
}
type Bar struct {
	Foo Foo
	C   string
}

func TestGetExpectTypeValue(t *testing.T) {
	a := Bar{Foo: Foo{A: 1}}
	b := Foo{}
	v := GetFieldValue(&a, &b)
	if v {
		log.Info(b)
	}
}
