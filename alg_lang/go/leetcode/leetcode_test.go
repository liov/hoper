package leetcode

import (
	"fmt"
	"testing"
)

func TestIsValid(t *testing.T) {
	println(isValid("()()()"))
}

func TestCompareVersion(t *testing.T) {
	println(compareVersion2("11", "10"))
}

func TestLongestPalindrome(t *testing.T) {
	println(longestPalindrome("ac"))
}

func TestConvert(t *testing.T) {
	println(convert("PAYPALISHIRING", 3))
}

func TestReverseKGroup(t *testing.T) {
	list := NewList([]int{1, 2})
	fmt.Println(reverseKGroup(list, 2))
}

func TestMyAtoi(t *testing.T) {
	fmt.Println(myAtoi("00000-42a1234"))
}

func TestSearch2(t *testing.T) {
	fmt.Println(search2([]int{8, 1, 2, 3, 4, 5, 6, 7}, 6))
}
