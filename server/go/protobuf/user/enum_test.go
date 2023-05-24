package user

import (
	"fmt"
	"testing"
)

func Test_Enum(t *testing.T) {
	var a interface{}
	a = Role(0)
	fmt.Println(a)
	b := Role(1)
	fmt.Println(b.String())
	fmt.Println(b.OrigString())
}
