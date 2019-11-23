package main

/*
type nonEmptyInterface struct {
	// see ../runtime/iface.go:/Itab
	itab *struct {
		ityp *rtype // static interface type
		typ  *rtype // dynamic concrete type
		hash uint32 // copy of typ.hash
		_    [4]byte
		fun  [100000]unsafe.Pointer // method table
	}
	word unsafe.Pointer
}
type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
*/
type Foo interface {
	print()
}

