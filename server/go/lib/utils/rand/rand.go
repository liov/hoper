package rand

import "math/rand"

func Intn(min, max int) int {
	return rand.Intn(max-min) + min
}
