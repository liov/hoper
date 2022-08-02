package tool

import "math"

var table string = "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"
var s = [6]int{11, 10, 3, 8, 4, 6}
var xor = 177451812
var add = 8728348608
var tr map[string]int

// source code: https://blog.csdn.net/dotastar00/article/details/108805779
func Bv2av(x string) int64 {
	tr = make(map[string]int)
	for i := 0; i < 58; i++ {
		tr[string(table[i])] = i
	}
	r := 0
	for i := 0; i < 6; i++ {
		r += tr[string(x[s[i]])] * int(math.Pow(float64(58), float64(i)))
	}
	return int64((r - add) ^ xor)
}

func Av2bv(x int) string {
	x = (x ^ xor) + add
	r := []byte("BV1  4 1 7  ")
	for i := 0; i < 6; i++ {
		r[s[i]] = table[int(math.Floor(float64(x)/float64(int(math.Pow(58, float64(i)))%58)))]
	}

	return string(r)
}
