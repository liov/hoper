package model

import (
	"fmt"
	"testing"
)

func Test_Enum(t *testing.T) {
	var a interface{}
	a = Role(0)
	fmt.Println(a)
}
