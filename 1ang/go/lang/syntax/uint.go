package main

import (
	"encoding/binary"
	"fmt"
	"github.com/hopeio/cherry/utils/number"
	"strconv"
)

// int类型之间强转，短位转高位负数前面补1正数前面补0，高位转短位截取
func main() {
	var i int32 = 255
	number.ViewBin(i)
	number.ViewBin(uint64(i))
	number.ViewBin(int64(i))
	number.ViewBin(int8(i)) //-1
	fmt.Println(byte(i))
	fmt.Println(string(i))
	var ii int64 = 255
	var b = make([]byte, 8, 8)
	binary.BigEndian.PutUint64(b, uint64(ii))
	fmt.Println(string(b))
	var arr [8]byte
	PutUint64(&arr, ii)
	fmt.Println(arr)
	fmt.Printf("%032b\n", i)
	fmt.Printf("%016o\n", i)
	fmt.Printf("%08x\n", i)
	fmt.Println(strconv.FormatInt(ii, 16))
}

func PutUint64(b *[8]byte, v int64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}
