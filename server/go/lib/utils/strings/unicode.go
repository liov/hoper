package stringsi

import (
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	runei "github.com/liov/hoper/server/go/lib/utils/strings/rune"
)

// [。；，：“”（）、？《》]
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

// unquote converts a quoted JSON string literal s into an actual string t.
// The rules are different than for Go, so cannot use strconv.Unquote.
func Unquote(s []byte) (t string, ok bool) {
	s, ok = unquoteBytes(s)
	t = ToString(s)
	return
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}

func ConvertUnicode(s []byte) string {
	if len(s) < 6 {
		return ToString(s)
	}
	b := make([]byte, len(s)+2*utf8.UTFMax)
	begin, bbegin := 0, 0
	for i := 0; i+6 < len(s); {
		if s[i] == '\\' && s[i+1] == 'u' {
			bbegin += copy(b[bbegin:], s[begin:i])
			rr := getu4(s[i:])
			if rr < 0 {
				return ToString(s)
			}
			i += 6
			if utf16.IsSurrogate(rr) {
				rr1 := getu4(s[i:])
				if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
					// A valid pair; consume.
					i += 6
					bbegin += utf8.EncodeRune(b[bbegin:], dec)
					break
				}
				// Invalid surrogate; fall back to replacement rune.
				rr = unicode.ReplacementChar
			}
			begin = i
			bbegin += utf8.EncodeRune(b[bbegin:], rr)
		} else {
			i++
		}
	}
	bbegin += copy(b[bbegin:], s[begin:])
	return ToString(b[:bbegin])
}
