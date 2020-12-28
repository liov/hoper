package aop

import (
	"log"
	"reflect"
	"unsafe"
)

func Aop(before, fs, after []func()) {
	for _, f := range before {
		f()
	}
	for _, f := range fs {
		f()
	}
	for _, f := range after {
		f()
	}
}

func Invoke(aop func(), self interface{}) {
	v2 := reflect.ValueOf(self).Elem()
	if v2.Kind() != reflect.Func {
		panic("错误的类型")
	}
	log.Println(v2.Pointer())
	code := reflect.ValueOf(v2).FieldByName("ptr").Pointer()
	ptr := *(*unsafe.Pointer)(*(*unsafe.Pointer)(unsafe.Pointer(code)))
	oldFuncVal := reflect.MakeFunc(v2.Type(), nil)
	funcValuePtr := reflect.ValueOf(oldFuncVal).FieldByName("ptr").Pointer()
	funcPtr := (*Func)(unsafe.Pointer(funcValuePtr))
	funcPtr.codePtr = uintptr(ptr)
	newFuncVal := reflect.MakeFunc(v2.Type(), func(in []reflect.Value) []reflect.Value {
		aop()
		return oldFuncVal.Call(in)
	})
	v2.Set(newFuncVal)

	log.Println(v2.Pointer())
}

type Func struct {
	codePtr uintptr
}
