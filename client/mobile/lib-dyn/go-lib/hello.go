package main

import "C"

//export goprint
func goprint(s *C.char) *C.char {
	return (*C.char)(C.CString("From DLL: " + C.GoString(s)))
}

func main() {
	// Need a main function to make CGO compile package as C shared library
}
