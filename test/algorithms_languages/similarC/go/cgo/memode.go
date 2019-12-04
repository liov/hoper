package main

/*
#include <stdlib.h>

void* makeslice(size_t memsize) {
    return malloc(memsize);
}
*/
import "C"
import "unsafe"

//创建大于2GB的切片
func makeByteSlize(n int) []byte {
	p := C.makeslice(C.size_t(n))
	return ((*[1 << 32]byte)(p))[0:]
}

func freeByteSlice(p []byte) {
	C.free(unsafe.Pointer(&p[0]))
}

func main() {
	s := makeByteSlize(1<<32 + 1)
	println(1<<32, len(s))
	s[len(s)-1] = 255
	println(s[len(s)-1])
	freeByteSlice(s)
}
