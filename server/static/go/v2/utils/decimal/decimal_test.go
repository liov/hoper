package decimal

import (
	"fmt"
	"testing"
)

func Test_Dec(t *testing.T) {
	fmt.Printf("%#v", Decimal{exponent: 6})
}

func Test_Decimal(t *testing.T) {
	a, _ := New("2.005")
	b, _ := New(" 0.1000000000000000055511151231257827021181583404541015625")
	fmt.Println(a, b)
}

func Test_Float(t *testing.T) {
	var a = 0.1000000000000000055511151231257827021181583404541015625
	fmt.Println(a)
}
