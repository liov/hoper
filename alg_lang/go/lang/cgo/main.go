//go:build go1.12
// +build go1.12

package main

//generate ${SRCDIR} 不能用,值是""
//go:generate cmd /C echo ${SRCDIR}
//go:generate gcc -c -o /demo/dll/hello.o hello.c
//go:generate ar rcs /demo/dll/libhello.a ${SRCDIR}/demo/dll/hello.o

//

//go run test/cgo

/*
#cgo  CFLAGS: -I${SRCDIR}/../../../../../tool/mingw/mingw64/include -I./include
#cgo  LDFLAGS: -L${SRCDIR}/../../../../dll -lhello
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include "hello.h"

void SayHelloInner(_GoString_ s);

static void SayHello(const char* s) {
    puts(s);
}

struct A {
	int i;
 	float f;
    int type;// type 是 Go 语言的关键字
	float _type; // 将屏蔽CGO对 type 成员的访问
  	int   size: 10; // 位字段无法访问,需要通过在C语言中定义辅助函数来完成
    float arr[];    // 零长的数组也无法访问,但其中零长的数组成员所在位置的偏移量依然可以通过unsafe.Offsetof(a.arr)来访问
};

union B1 {
    int i;
    float f;
};

union B2 {
    int8_t i8;
    int64_t i64;
};

union B {
    int i;
    float f;
};

enum C {
    ONE,
    TWO,
};

//char *s1 = "Hello"; //不能使static的
char arr[10];

static int add(int a, int b) {
    return a+b;
}
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	C.SayHello(C.CString("Hello, World\n"))
	C.SayHelloExternal(C.CString("Hello, World\n"))
	C.SayHelloInner("Hello, World\n")
	var a C.struct_A
	fmt.Println(a.i)
	fmt.Println(a.f)
	fmt.Println(a._type)

	var b1 C.union_B1
	fmt.Printf("%T\n", b1) // [4]uint8

	var b2 C.union_B2
	fmt.Printf("%T\n", b2) // [8]uint8

	var b C.union_B
	fmt.Println("b.i:", *(*C.int)(unsafe.Pointer(&b)))
	fmt.Println("b.f:", *(*C.float)(unsafe.Pointer(&b)))

	var c C.enum_C = C.TWO
	fmt.Println(c)
	fmt.Println(C.ONE)
	fmt.Println(C.TWO)

	// 通过 reflect.SliceHeader 转换
	var arr0 []byte
	var arr0Hdr = (*reflect.SliceHeader)(unsafe.Pointer(&arr0))
	arr0Hdr.Data = uintptr(unsafe.Pointer(&C.arr[0]))
	arr0Hdr.Len = 10
	arr0Hdr.Cap = 10

	// 通过切片语法转换
	arr1 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:10:10]

	/*	var s0 string
		var s0Hdr = (*reflect.StringHeader)(unsafe.Pointer(&s0))
		s0Hdr.Result = uintptr(unsafe.Pointer(C.s1))
		s0Hdr.Len = int(C.strlen(C.s1))

		sLen := int(C.strlen(C.s1))
		s1 := string((*[31]byte)(unsafe.Pointer(&C.s1[0]))[:sLen:sLen])*/
	fmt.Println(arr1, C.add(1, 1))

	CallC()
}

func CallC() {
	C.SayHello(C.CString("Hello, World\n"))
}

//export SayHelloInner
func SayHelloInner(s string) {
	fmt.Print(s)
}

//export helloString
func helloString(s string) {}

//export helloSlice
func helloSlice(s []byte) {}
