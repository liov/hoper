package _type

import (
	"unsafe"
)

type Rtype struct {
	Size       uintptr
	ptrdata    uintptr  // number of bytes in the type that can contain pointers
	hash       uint32   // hash of type; avoids computation in hash tables
	tflag      tflag    // extra type information flags
	align      uint8    // alignment of variable with this type
	fieldAlign uint8    // alignment of struct field with this type
	Kind       uint8    // enumeration for C
	alg        *typeAlg // algorithm table
	gcdata     *byte    // garbage collection data
	str        nameOff  // string form
	ptrToThis  typeOff  // type for pointer to this type, may be zero
}

type tflag uint8
type nameOff int32 // offset to a name
type typeOff int32 // offset to an *Rtype
type textOff int32 // offset from top of text section
// a copy of runtime.typeAlg
type typeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type method struct {
	name nameOff // name of method
	mtyp typeOff // method type (without receiver)
	ifn  textOff // fn used in interface call (one-word receiver)
	tfn  textOff // fn used for normal method call
}

// uncommonType is present only for defined types or types with methods
// (if T is a defined type, the uncommonTypes for T and *T have methods).
// Using a pointer to this struct reduces the overall Size required
// to describe a non-defined type with no methods.
type uncommonType struct {
	pkgPath nameOff // import path; empty for built-in types like int, string
	mcount  uint16  // number of methods
	xcount  uint16  // number of exported methods
	moff    uint32  // offset from this uncommontype to [mcount]method
	_       uint32  // unused
}

// ChanDir represents a channel type's direction.
type ChanDir int

type arrayType struct {
	Rtype
	elem  *Rtype // array element type
	slice *Rtype // slice type
	len   uintptr
}

// chanType represents a channel type.
type chanType struct {
	Rtype
	elem *Rtype  // channel element type
	dir  uintptr // channel direction (ChanDir)
}

type funcType struct {
	Rtype
	inCount  uint16
	outCount uint16 // top bit is set if last input parameter is ...
}

// imethod represents a method on an interface type
type imethod struct {
	name nameOff // name of method
	typ  typeOff // .(*FuncType) underneath
}

// interfaceType represents an interface type.
type interfaceType struct {
	Rtype
	pkgPath name      // import path
	methods []imethod // sorted by hash
}

// mapType represents a map type.
type mapType struct {
	Rtype
	key        *Rtype // map key type
	elem       *Rtype // map element (value) type
	bucket     *Rtype // internal bucket structure
	keySize    uint8  // Size of key slot
	valueSize  uint8  // Size of value slot
	bucketSize uint16 // Size of bucket
	flags      uint32
}

// ptrType represents a pointer type.
type ptrType struct {
	Rtype
	elem *Rtype // pointer element (pointed at) type
}

// sliceType represents a slice type.
type sliceType struct {
	Rtype
	elem *Rtype // slice element type
}

// Struct field
type structField struct {
	name        name    // name is always non-empty
	typ         *Rtype  // type of field
	offsetEmbed uintptr // byte offset of field<<1 | isEmbedded
}

type structType struct {
	Rtype
	pkgPath name
	fields  []structField // sorted by offset
}

type name struct {
	bytes *byte
}

type EmptyInterface struct {
	Typ  *Rtype
	Word unsafe.Pointer
}
