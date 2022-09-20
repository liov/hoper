package main

import "fmt"

var a = 15
var b = 16

func main() {
	fmt.Println(a&1, b&1)
	fmt.Println(a|b, a|1)
	fmt.Println(a^b, a^a, a^b^b)
	fmt.Println(^a, ^b, 0^a)
	fmt.Println(a&^b, (a^b)&b) // 标志位操作 &^,清除标记位
	fmt.Println(a<<2, b<<2)
	fmt.Println(a>>2, b>>2)
	var a uint8 = 0x82
	var b uint8 = 0x02
	fmt.Printf("%08b [A]\n", a)
	fmt.Printf("%08b [B]\n", b)

	fmt.Printf("%08b (NOT B)\n", ^b)
	fmt.Printf("%08b ^ %08b = %08b [B XOR 0xff]\n", b, 0xff, b^0xff)

	fmt.Printf("%08b ^ %08b = %08b [A XOR B]\n", a, b, a^b)
	fmt.Printf("%08b & %08b = %08b [A AND B]\n", a, b, a&b)
	fmt.Printf("%08b &^%08b = %08b [A 'AND NOT' B]\n", a, b, a&^b)
	fmt.Printf("%08b&(^%08b)= %08b [A AND (NOT B)]\n", a, b, a&(^b))
	fmt.Printf("0x2 & 0x2 + 0x4 -> %#x\n", 0x2&0x2+0x4) // & 优先 +
	//prints: 0x2 & 0x2 + 0x4 -> 0x6
	//Go:    (0x2 & 0x2) + 0x4
	//C++:    0x2 & (0x2 + 0x4) -> 0x2

	fmt.Printf("0x2 + 0x2 << 0x1 -> %#x\n", 0x2+0x2<<0x1) // << 优先 +
	//prints: 0x2 + 0x2 << 0x1 -> 0x6
	//Go:     0x2 + (0x2 << 0x1)
	//C++:   (0x2 + 0x2) << 0x1 -> 0x8

	fmt.Printf("0xf | 0x2 ^ 0x2 -> %#x\n", 0xf|0x2^0x2) // | 优先 ^
	//prints: 0xf | 0x2 ^ 0x2 -> 0xd
	//Go:    (0xf | 0x2) ^ 0x2
	//C++:    0xf | (0x2 ^ 0x2) -> 0xf

	var x uint8 = 1<<1 | 1<<5
	var y uint8 = 1<<1 | 1<<2

	fmt.Printf("%08b\n", x) // "00100010", the set {1, 5}
	fmt.Printf("%08b\n", y) // "00000110", the set {1, 2}

	fmt.Printf("%08b\n", x&y)  // "00000010", the intersection {1}
	fmt.Printf("%08b\n", x|y)  // "00100110", the union {1, 2, 5}
	fmt.Printf("%08b\n", x^y)  // "00100100", the symmetric difference {2, 5}
	fmt.Printf("%08b\n", x&^y) // "00100000", the difference {5}

	for i := uint(0); i < 8; i++ {
		if x&(1<<i) != 0 { // membership test
			fmt.Println(i) // "1", "5"
		}
	}

	fmt.Printf("%08b\n", x<<1) // "01000100", the set {2, 6}
	fmt.Printf("%08b\n", x>>1) // "00010001", the set {0, 4}
	//与非，清除，清除a中对位b中为1的位
	fmt.Printf("%08b &^%08b = %08b [A 'AND NOT' B]\n", 0x06, 0xff, 0x06&^0xff)
	fmt.Printf("%08b&(^%08b)= %08b [A AND (NOT B)]\n", 0x06, 0xff, 0x06&(^0xff))
	fmt.Printf("%08b ^ %08b = %08b [A XOR B]\n", 0x06, 0xff, 0x06^0xff)
}
