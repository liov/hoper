package stringsi

import (
	"unicode"

	runei "github.com/liov/hoper/server/go/lib/utils/strings/rune"
)

//[。；，：“”（）、？《》]
var HanPunctuation = []rune{
	'\u3002', '\uff1b', '\uff0c', '\uff1a', '\u201c', '\u201d', '\uff08', '\uff09', '\u3001', '\uff1f', '\u300a', '\u300b',
}

func HasHan(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || runei.In(r, HanPunctuation) {
			return true
		}
	}
	return false
}
