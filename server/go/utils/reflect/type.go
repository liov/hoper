package reflecti

import (
	"reflect"
	"unsafe"
)

type EmptyInterface struct {
	Typ  *Rtype
	Word unsafe.Pointer
}

// NonEmptyInterface is the header for an interface value with methods.
type NonEmptyInterface struct {
	// see ../runtime/iface.go:/Itab
	Itab *struct {
		Ityp *Rtype // static interface type
		Typ  *Rtype // dynamic concrete type
		Hash uint32 // copy of typ.hash
		_    [4]byte
		Fun  [100000]unsafe.Pointer // method table
	}
	Word unsafe.Pointer
}

type Iface struct {
	Itab *Itab
	Data unsafe.Pointer
}

type Eface struct {
	Type *Type
	Data unsafe.Pointer
}

type Itab struct {
	Inter  *InterfaceType
	Type   *Type
	Link   *Itab
	Hash   uint32 // copy of _type.hash. Used for type switches.
	Bad    bool   // type does not implement interface
	Inhash bool   // has this Itab been added to hash?
	Unused [2]byte
	Fun    [1]uintptr // variable sized
}

// interfaceType represents an interface type.
type InterfaceType struct {
	Rtype
	PkgPath Name      // import path
	Methods []Imethod // sorted by hash
}

type Type struct {
	Size       uintptr
	Ptrdata    uintptr // size of memory prefix holding all pointers
	Hash       uint32
	Tflag      Tflag
	Align      uint8
	FieldAlign uint8
	Kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	Gcdata    *byte
	Str       NameOff
	PtrToThis TypeOff
}

// imethod represents a method on an interface type
type Imethod struct {
	Name NameOff // name of method
	Typ  TypeOff // .(*FuncType) underneath
}

type Tflag uint8
type NameOff int32 // offset to a name
type TypeOff int32 // offset to an *Rtype
type textOff int32 // offset from top of text section

// Rtype must be kept in sync with ../runtime/type.go:/^type._type.
type Rtype struct {
	Size       uintptr
	Ptrdata    uintptr // number of bytes in the type that can contain pointers
	Hash       uint32  // hash of type; avoids computation in hash tables
	tflag      Tflag   // extra type information flags
	Align      uint8   // alignment of variable with this type
	FieldAlign uint8   // alignment of struct field with this type
	Kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal     func(unsafe.Pointer, unsafe.Pointer) bool
	Gcdata    *byte   // garbage collection data
	Str       NameOff // string form
	PtrToThis TypeOff // type for pointer to this type, may be zero
}

func RtypeOff(section unsafe.Pointer, off int32) *Rtype {
	return (*Rtype)(add(section, uintptr(off), "sizeof(rtype) > 0"))
}

func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// arrayType represents a fixed array type.
type arrayType struct {
	Rtype
	Elem  *Rtype // array element type
	Slice *Rtype // Slice type
	Len   uintptr
}

// ChanType represents a channel type.
type ChanType struct {
	Rtype
	Elem *Rtype  // channel element type
	Dir  uintptr // channel direction (ChanDir)
}

// funcType represents a function type.
//
// A *Rtype for each in and out parameter is stored in an array that
// directly follows the funcType (and possibly its uncommonType). So
// a function type with one method, one input, and one output is:
//
//	struct {
//		funcType
//		uncommonType
//		[2]*Rtype    // [0] is in, [1] is out
//	}
type funcType struct {
	Rtype
	InCount  uint16
	OutCount uint16 // top bit is set if last input parameter is ...
}

// MapType represents a map type.
type MapType struct {
	Rtype
	Key    *Rtype // map key type
	Elem   *Rtype // map element (value) type
	Bucket *Rtype // internal bucket structure
	// function for hashing keys (ptr to key, seed) -> hash
	Hasher     func(unsafe.Pointer, uintptr) uintptr
	Keysize    uint8  // size of key slot
	Valuesize  uint8  // size of value slot
	Bucketsize uint16 // size of bucket
	Flags      uint32
}

// PtrType represents a pointer type.
type PtrType struct {
	Rtype
	Elem *Rtype // pointer element (pointed at) type
}

// SliceType represents a Slice type.
type SliceType struct {
	Rtype
	Elem *Rtype // Slice element type
}

// StructType represents a struct type.
type StructType struct {
	Rtype
	PkgPath Name
	Fields  []StructField // sorted by offset
}

// Struct field
type StructField struct {
	Name        Name    // name is always non-empty
	Typ         *Rtype  // type of field
	OffsetEmbed uintptr // byte offset of field<<1 | isEmbedded
}

type Name struct {
	Bytes *byte
}

type FuncID uint8

type Func struct {
	Entry   uintptr // start pc
	Nameoff int32   // function name

	Args        int32  // in/out args size
	Deferreturn uint32 // offset of start of a deferreturn call instruction from entry, if any.

	Pcsp      uint32
	Pcfile    uint32
	Pcln      uint32
	Npcdata   uint32
	CuOffset  uint32  // runtime.cutab offset of this function's CU
	FuncID    FuncID  // set for certain special runtime functions
	_         [2]byte // pad
	Nfuncdata uint8   // must be last
}

func Hash(v interface{}) uint32 {
	ia := *(*Iface)(unsafe.Pointer(&v))
	return ia.Itab.Hash
}

type Flag uintptr

type Value struct {
	Typ *Rtype
	Ptr unsafe.Pointer
	Flag
}

func RuntimeTypeID(t reflect.Type) uintptr {
	return uintptrElem(uintptr(unsafe.Pointer(&t)) + PtrOffset)
}

func DereferenceType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
