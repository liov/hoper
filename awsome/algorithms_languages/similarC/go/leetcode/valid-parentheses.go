package leetcode

func isValid(s string) bool {
	var balance int
	for _, c := range s {
		if c == '(' {
			balance++
		} else {
			balance--
		}
		if balance < 0 {
			return false
		}
		if balance == 0 {
			return true
		}
	}
	return false
}
