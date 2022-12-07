package binary

import (
	reflecti "github.com/liov/hoper/server/go/lib/utils/reflect"
	"reflect"
	"unsafe"
)

/*
 *这种序列化序列化的是指针，一旦结构体包含指针，是不能从[]byte里还原结构体的，只能做临时的转换，基本没什么用，序列化和
 *反序列化必须成对出现，而且go的GC偏移回收的话，有可能也GG
 */

func getSize(t interface{}) int {
	size := reflect.TypeOf(t).Elem().Size()
	return (int)(size)
}

func StructToBytes(s interface{}) []byte {
	sizeOfStruct := getSize(s)
	var x reflect.SliceHeader
	x.Len = sizeOfStruct
	x.Cap = sizeOfStruct
	x.Data = uintptr((*reflecti.Eface)(unsafe.Pointer(&s)).Data)
	return *(*[]byte)(unsafe.Pointer(&x))
}

func BytesToMyStruct(b []byte) unsafe.Pointer {
	return unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	)
}
