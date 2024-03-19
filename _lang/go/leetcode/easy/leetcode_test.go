package easy

import (
	"fmt"
	"testing"
)

func TestIsValid(t *testing.T) {
	println(isValid("()"))
}

func TestLargestSumAfterKNegations(t *testing.T) {
	fmt.Println(largestSumAfterKNegations([]int{2, -3, -1, 5, -4}, 2))
}
