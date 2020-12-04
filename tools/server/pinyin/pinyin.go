package py

import (
	"unicode/utf8"

	"github.com/mozillazg/go-pinyin"
)

var FirstFallBack = func(r rune, a pinyin.Args) []string {
	if r < utf8.RuneSelf && ('A' <= byte(r) && byte(r) <= 'Z') {
		r += 'a' - 'A'
	}
	return []string{string(r)}
}

// Pinyin 汉字转拼音，支持多音字模式.
func FistLetter(s string) string {
	a := pinyin.Args{Style: pinyin.FirstLetter, Heteronym: pinyin.Heteronym, Separator: pinyin.Separator, Fallback: FirstFallBack}

	for _, r := range s {
		py := pinyin.SinglePinyin(r, a)
		if len(py) > 0 {
			return py[0]
		}
	}
	return "-"
}
