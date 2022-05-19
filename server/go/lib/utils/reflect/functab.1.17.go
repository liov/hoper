//go:build !go1.18

package reflecti

type functab struct {
	entry   uintptr
	funcoff uintptr
}
