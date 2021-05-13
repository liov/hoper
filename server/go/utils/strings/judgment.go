package stringsi

// Is c an ASCII lower-case letter?
func IsASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func IsASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func In(s string, ss []string) bool {
	for _, rn := range ss {
		if s == rn {
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
