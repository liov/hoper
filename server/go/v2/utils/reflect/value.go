package reflecti

import (
	"reflect"
	"unsafe"
)

const (
	FlagKindWidth        = 5 // there are 27 kinds
	FlagKindMask    Flag = 1<<FlagKindWidth - 1
	FlagStickyRO    Flag = 1 << 5
	FlagEmbedRO     Flag = 1 << 6
	FlagIndir       Flag = 1 << 7
	FlagAddr        Flag = 1 << 8
	FlagMethod      Flag = 1 << 9
	FlagMethodShift      = 10
	FlagRO          Flag = FlagStickyRO | FlagEmbedRO
)

var (
	e         = EmptyInterface{Typ: new(Rtype)}
	PtrOffset = func() uintptr {
		return unsafe.Offsetof(e.Word)
	}()
	KindOffset = func() uintptr {
		return unsafe.Offsetof(e.Typ.Kind)
	}()
	ElemOffset = func() uintptr {
		return unsafe.Offsetof(new(PtrType).Elem)
	}()
	SliceDataOffset = func() uintptr {
		return unsafe.Offsetof(new(reflect.SliceHeader).Data)
	}()
)

// DereferenceValue dereference and unpack interface,
// get the underlying non-pointer and non-interface value.
func DereferenceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DereferencePtrValue returns the underlying non-pointer type value.
func DereferencePtrValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// DereferenceInterfaceValue returns the value of the underlying type that implements the interface v.
func DereferenceInterfaceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

//go:nocheckptr
func ValueOf(v interface{}) Value {
	stdValue := reflect.ValueOf(v)
	return *(*Value)(unsafe.Pointer(&stdValue))
}

//go:nocheckptr
func ConvertValue(v reflect.Value) Value {
	return *(*Value)(unsafe.Pointer(&v))
}

//go:nocheckptr
func getFlag(typPtr uintptr) Flag {
	if unsafe.Pointer(typPtr) == nil {
		return 0
	}
	return *(*Flag)(unsafe.Pointer(typPtr + KindOffset))
}

//go:nocheckptr
func pointerElem(p unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(p)
}

// Pointer gets the pointer of i.
// NOTE:
//  *T and T, gets diffrent pointer
//go:nocheckptr
func (v Value) Pointer() uintptr {
	switch v.Kind() {
	case reflect.Invalid:
		return 0
	case reflect.Slice:
		return uintptrElem(uintptr(v.Ptr)) + SliceDataOffset
	default:
		return uintptr(v.Ptr)
	}
}

// Kind gets the reflect.Kind fastly.
func (v Value) Kind() reflect.Kind {
	return reflect.Kind(v.Flag & FlagKindMask)
}

//go:nocheckptr
func uintptrElem(ptr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr))
}
