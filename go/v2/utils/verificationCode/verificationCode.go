package verificationCode

import (
	"math/rand"
	"time"
)

var code = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K',
	'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func Generate() string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(code),  func(i, j int) {
		code[i], code[j] = code[j], code[i]
	})
	return string(code[:4])
}
