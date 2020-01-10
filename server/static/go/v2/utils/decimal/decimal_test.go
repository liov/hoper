package decimal

import (
	"fmt"
	"testing"
)

func Test_Dec(t *testing.T) {
	fmt.Printf("%#v", dec{exponent: 6})
}

func Test_Decimal(t *testing.T) {
	a, _ := New(2, "005")
	b, _ := New(1, "5")
	fmt.Println(a, b)
}
