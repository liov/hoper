package main

/*
#include <stdint.h>

int64_t myadd(int64_t a, int64_t b) {
    return a+b;
}
*/
import "C"
import (
	"unsafe"

	asmpkg "test/cgo/asm"
)

func main() {

	println(asmpkg.AsmCallCAdd(
		uintptr(unsafe.Pointer(C.myadd)),
		123, 456,
	))
	asmpkg.SyscallWrite_Darwin(1, "hello syscall!\n")

}
