package reflecti

import (
	"log"
	"testing"
)

func TestGetFunc(t *testing.T) {
	var fn func(sep string) uint32
	err := GetFunc(&fn, "strings.hashStr")
	if err != nil {
		log.Println(err)
	}
	log.Println(fn("test"))
}
