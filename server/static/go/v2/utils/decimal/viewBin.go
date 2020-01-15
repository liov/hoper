package decimal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

/*Sizeof函数返回的大小只包括数据结构中固定的部分，例如字符串对应结构体中的指针和字符串长度部分，但是并不包含指针指向的字符串的内容。Go语言中非聚合类型通常有一个固定的大小，尽管在不同工具链下生成的实际大小可能会有所不同。考虑到可移植性，引用类型或包含引用类型的大小在32位平台上是4个字节，在64位平台上是8个字节。

计算机在加载和保存数据时，如果内存地址合理地对齐的将会更有效率。例如2字节大小的int16类型的变量地址应该是偶数，一个4字节大小的rune类型变量的地址应该是4的倍数，一个8字节大小的float64、uint64或64-bit指针类型变量的地址应该是8字节对齐的。但是对于再大的地址对齐倍数则是不需要的，即使是complex128等较大的数据类型最多也只是8字节对齐。

由于地址对齐这个因素，一个聚合类型（结构体或数组）的大小至少是所有字段或元素大小的总和，或者更大因为可能存在内存空洞。*/
/*1字节（byte）=8位（bit）
在16位系统中，1字（word）=2字节（byte）=16位（bit）
在32位系统中，1字（word）=4字节（byte）=32位（bit）
在64位系统中，1字（word）=8字节（byte）=64位（bit）*/
func ViewBin(v interface{}) {
	vv := reflect.ValueOf(v)
	var binary string
	switch v.(type) {
	case int, int8, int16, int32, int64:
		if vv.Int() < 0 {
			f := fmt.Sprintf("%064b", uint64(vv.Int()))
			binary = f[len(f)-int(vv.Type().Size())*8:]
		} else {
			binary = fmt.Sprintf("%0"+strconv.Itoa(int(vv.Type().Size())*8)+"b", v)
		}
	case uint, uint8, uint16, uint32, uint64:
		binary = fmt.Sprintf("%0"+strconv.Itoa(int(vv.Type().Size())*8)+"b", v)
	case float32, float64:
		f := vv.Float()
		ViewBin(*(*int64)(unsafe.Pointer(&f)))
		return
	}
	var out []string
	for i := 0; i < 8; i++ {
		out = append(out, binary[i*8:(i+1)*8])
	}
	fmt.Println(strings.Join(out, " "), " ", v)
}
