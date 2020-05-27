package strings2

import (
	"log"
	"testing"
)

func TestFormatLen(t *testing.T) {
	s := "post"
	log.Println(FormatLen(s, 10), "test")
	s = "AutoCommit"
	log.Println(ConvertToSnackCase(s))
}
