package stringsi

import (
	"strings"
	"unsafe"
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

func ConvertToSnackCase(in string) string {
	const (
		lower = false
		upper = true
	)

	if in == "" {
		return ""
	}
	in = strings.TrimSpace(in)
	var (
		buf                                      = new(strings.Builder)
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range in[:len(in)-1] {
		nextCase = in[i+1] >= 'A' && in[i+1] <= 'Z'
		nextNumber = in[i+1] >= '0' && in[i+1] <= '9'

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if in[i-1] != '_' && in[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(in)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(in[len(in)-1])

	s := strings.ToLower(buf.String())
	return s
}

func ConvertToCamelCase(in string) string {
	if in == "" {
		return ""
	}
	in = strings.TrimSpace(in)
	buf := new(strings.Builder)
	nextCaseUp := false
	skip := false
	buf.WriteByte(LowerCase(in[0]))
	for i, _ := range in[:len(in)-1] {
		if !skip {
			if nextCaseUp {
				buf.WriteByte(UpperCase(in[i]))
				nextCaseUp = false
			} else {
				buf.WriteByte(in[i])
			}
		} else {
			nextCaseUp = true
			continue
		}
		skip = in[i+1] == '-' || in[i+1] == '_' || in[i+1] == '.'
	}
	return buf.String()
}

// 仅首位小写（更符合接口的规范）
func LowerFirst(t string) string {
	b := []byte(t)
	if 'A' <= b[0] && b[0] <= 'Z' {
		b[0] += 'a' - 'A'
	}
	return string(b)
}

func LowerCase(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	return c
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
	needCoby := false
	last := false
	for i, r := range s {
		if in(r, old) {
			if needCoby {
				w += copy(t[w:], s[start:i])
				needCoby = false
			}
			last = true
			continue
		}
		needCoby = true
		if last {
			start = i
			last = false
		}
	}
	if needCoby {
		w += copy(t[w:], s[start:])
	}
	return string(t[0:w])
}

func in(r rune, bytes []rune) bool {
	for i := range bytes {
		if bytes[i] == r {
			return true
		}
	}
	return false
}

// 有一个匹配成功就返回true
func HasPrefixes(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
			return true
		}
	}
	return false
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

// And now lots of helper functions.

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
