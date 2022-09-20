package main

import "unicode/utf8"

func main() {
	s := "我爱中国"
	for i, r := range s {
		println(i, r, utf8.RuneLen(r))
	}
}
