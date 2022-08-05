package client

import "math/rand"

const (
	UserAgent1 = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0"
)

var userAgent = []string{UserAgent1}

func GetRandUserAgent() string {
	return userAgent[rand.Intn(len(userAgent))]
}
