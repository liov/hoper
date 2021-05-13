package binary

import (
	"unsafe"
)

func IntFromBinary(b []byte) int {
	return *(*int)(unsafe.Pointer(uintptr(*(*int)(unsafe.Pointer(&b)))))
}

func IntToBinary(i int) []byte {
	b := make([]byte, 8, 8)
	*(*int)(unsafe.Pointer(uintptr(*(*int)(unsafe.Pointer(&b))))) = i
	return b
}

// 比标准库慢很多,10倍左右，string和bytes互转只是节省复制内存，unsafe操作有很多检测
// binary.LittleEndian.Uint64(b)
func UIntFromBinary(b []byte) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(*(*int)(unsafe.Pointer(&b)))))
}

// binary.LittleEndian.PutUint64(b)
func UIntToBinary(i uint64) []byte {
	b := make([]byte, 8, 8)
	*(*uint64)(unsafe.Pointer(uintptr(*(*int)(unsafe.Pointer(&b))))) = i
	return b
}
