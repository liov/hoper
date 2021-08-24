package reflecti

import (
	"log"
	"reflect"
	"testing"
)

type Bar1 struct {
	Field1 int
	Field2 string `mock:"example:'1',type:'\\w'"`
}

func TestTag(t *testing.T) {
	var bar Bar1
	typ := reflect.TypeOf(bar)
	log.Println(GetCustomizeTag(typ.Field(1).Tag.Get("mock"), "example"))
}
