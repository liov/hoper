package h_reflect

import (
	"testing"

	"github.com/liov/hoper/go/v2/utils/log"
)

func TestGetExpectTypeValue(t *testing.T) {
	type Foo struct {
		A int
		B string
	}
	type Bar struct {
		Foo Foo
		C string
	}
	a:= Bar{Foo:Foo{A:1}}
	b:= Foo{}
	v:=GetExpectTypeValue(&a,&b)
	if v{
		log.Info(b)
	}
}
