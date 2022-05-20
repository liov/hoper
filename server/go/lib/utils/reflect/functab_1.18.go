//go:build go1.18 && !go1.19

package reflecti

type functab struct {
	entry   uint32
	funcoff uint32
}
