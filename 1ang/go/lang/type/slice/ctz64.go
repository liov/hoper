package main

import (
	"fmt"
	"math/bits"
)

func main() {
	fmt.Println(bits.TrailingZeros64(512))
	fmt.Println(MulUintptr(12, 13))
}

const MaxUintptr = ^uintptr(0)
const PtrSize = 4 << (^uintptr(0) >> 63)

// MulUintptr returns a * b and whether the multiplication overflowed.
// On supported platforms this is an intrinsic lowered by the compiler.
func MulUintptr(a, b uintptr) (uintptr, bool) {
	if a|b < 1<<(4*PtrSize) || a == 0 {
		return a * b, false
	}
	overflow := b > MaxUintptr/a
	return a * b, overflow
}
