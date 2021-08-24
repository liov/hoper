package number

import (
	"math"
	"time"
)

const magicNumber = 0xf1234fff

//一个数异或同一个数两次还是这个数...

func GenKey() int64 {
	return time.Now().Unix() ^ magicNumber
}

func Validation(key int64) float64 {
	return math.Abs(float64(key ^ magicNumber - time.Now().Unix()))
}
