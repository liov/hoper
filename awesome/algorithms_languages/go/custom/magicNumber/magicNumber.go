package main

import (
	"log"
	"math"
	"time"
)

const magicNumber = 0xf1234fff

//一个数异或同一个数两次还是这个数...
func main() {
	var userId = time.Now().Unix() ^ magicNumber
	log.Println(userId)
	var validation = func(key int64) float64 {
		return math.Abs(float64(key ^ magicNumber - time.Now().Unix()))
	}
	log.Println(validation(userId))
	log.Println(validation(1))
}
