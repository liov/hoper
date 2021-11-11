package leetcode

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func abs(a int) int {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}

func reverse(bytes []byte) []byte {
	ret := make([]byte, len(bytes))
	l := len(bytes) - 1
	for i := 0; i < len(bytes); i++ {
		ret[i] = bytes[l-i]
	}
	return ret
}

func reverseInPlace(bytes []byte) {
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1] = bytes[len(bytes)-1], bytes[i]
	}
}

func intSlice(i int, s []int) bool {
	for _, n := range s {
		if n == i {
			return true
		}
	}
	return false
}
