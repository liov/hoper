package stringsi

import (
	"bytes"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
	"unicode"
	"unsafe"

	runei "github.com/actliboy/hoper/server/go/lib/utils/strings/rune"
)

func FormatLen(s string, length int) string {
	if len(s) < length {
		return s + strings.Repeat(" ", length-len(s))
	}
	return s[:length]
}

func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func QuoteToBytes(s string) []byte {
	b := make([]byte, 0, len(s)+2)
	b = append(b, '"')
	b = append(b, []byte(s)...)
	b = append(b, '"')
	return b
}

func CamelToSnake(name string) string {
	var ret bytes.Buffer

	multipleUpper := false
	var lastUpper rune
	var beforeUpper rune

	for _, c := range name {
		// Non-lowercase character after uppercase is considered to be uppercase too.
		isUpper := (unicode.IsUpper(c) || (lastUpper != 0 && !unicode.IsLower(c)))

		if lastUpper != 0 {
			// Output a delimiter if last character was either the first uppercase character
			// in a row, or the last one in a row (e.g. 'S' in "HTTPServer").
			// Do not output a delimiter at the beginning of the name.

			firstInRow := !multipleUpper
			lastInRow := !isUpper

			if ret.Len() > 0 && (firstInRow || lastInRow) && beforeUpper != '_' {
				ret.WriteByte('_')
			}
			ret.WriteRune(unicode.ToLower(lastUpper))
		}

		// Buffer uppercase char, do not output it yet as a delimiter may be required if the
		// next character is lowercase.
		if isUpper {
			multipleUpper = (lastUpper != 0)
			lastUpper = c
			continue
		}

		ret.WriteRune(c)
		lastUpper = 0
		beforeUpper = c
		multipleUpper = false
	}

	if lastUpper != 0 {
		ret.WriteRune(unicode.ToLower(lastUpper))
	}
	return string(ret.Bytes())
}

func ConvertToCamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// 仅首位小写（更符合接口的规范）
func LowerFirst(t string) string {
	b := []byte(t)
	b[0] = LowerCase(b[0])
	return string(b)
}

func LowerCase(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	return c
}

func UpperCaseFirst(t string) string {
	b := []byte(t)
	b[0] = UpperCase(b[0])
	return string(b)
}

func UpperCase(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c
}

func ReplaceRuneEmpty(s string, old []rune) string {
	if len(old) == 0 {
		return s // avoid allocation
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s))
	w := 0
	start := 0
	needCopy := false
	last := false
	for i, r := range s {
		if runei.In(r, old) {
			if needCopy {
				w += copy(t[w:], s[start:i])
				needCopy = false
			}
			last = true
			continue
		}
		needCopy = true
		if last {
			start = i
			last = false
		}
	}
	if needCopy {
		w += copy(t[w:], s[start:])
	}
	return string(t[0:w])
}

func Camel(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && IsASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if IsASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if IsASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && IsASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// And now lots of helper functions.

func CamelCase[T ~string](s T) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Caller the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

func PrintComments(comments protogen.CommentSet, g *protogen.GeneratedFile) {
	for _, comment := range comments.LeadingDetached {
		g.P(comment)
	}
	g.P(comments.Leading)
	g.P(comments.LeadingDetached)
}

type NumLetterSlice[T any] ['z' - '0' + 1]T

// 原来数组支持这样用
func (n *NumLetterSlice[T]) Set(b byte, v T) {
	n[b-'0'] = v
}

func ReplaceBytesEmpty(s string, old []byte) string {
	if len(old) == 0 {
		return s // avoid allocation
	}
	tmpl := make([]bool, 255)

	for _, b := range old {
		tmpl[b] = true
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s))
	w := 0
	start := 0
	needCopy := false
	last := false
	for i, r := range s {
		if r < 256 && tmpl[r] {
			if needCopy {
				w += copy(t[w:], s[start:i])
				needCopy = false
			}
			last = true
			continue
		}
		needCopy = true
		if last {
			start = i
			last = false
		}
	}
	if needCopy {
		w += copy(t[w:], s[start:])
	}
	return string(t[0:w])
}
