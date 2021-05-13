package slices

import (
	"fmt"
	"testing"
)

func TestContains(t *testing.T) {
	val1 := []string{"a", "b", "c"}
	val2 := "a"
	val3 := "d"
	fmt.Println(StringContains(val1, val2))
	fmt.Println(StringContains(val1, val3))
}
