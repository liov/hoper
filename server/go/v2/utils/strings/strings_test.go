package stringsi

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

func TestReplaceRuneEmpty(t *testing.T) {
	s := "p我o爱s中t"
	log.Println(ReplaceRuneEmpty(s,[]rune{'o'}))
	log.Println(ReplaceRuneEmpty(s,[]rune{'o','s'}))
	log.Println(ReplaceRuneEmpty(s,[]rune{'o','t'}))
	log.Println(ReplaceRuneEmpty(s,[]rune{'中','t'}))
}