package stringsi

import (
	"unicode"
)

//[。；，：“”（）、？《》]
var HanPunctuation = []rune{
	'\u3002', '\uff1b', '\uff0c', '\uff1a', '\u201c', '\u201d', '\uff08', '\uff09', '\u3001', '\uff1f', '\u300a', '\u300b',
}

func HasHan(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || Contain(HanPunctuation, r) {
			return true
		}
	}
	return false
}

func Contain(runes []rune, r rune) bool {
	for _, rn := range runes {
		if r == rn {
			return true
		}
	}
	return false
}
