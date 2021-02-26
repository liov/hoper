package runei


func Contain(runes []rune, r rune) bool {
	for _, rn := range runes {
		if r == rn {
			return true
		}
	}
	return false
}

func In(r rune, bytes []rune) bool {
	for i := range bytes {
		if bytes[i] == r {
			return true
		}
	}
	return false
}
