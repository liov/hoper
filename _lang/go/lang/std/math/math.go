package main

import "log"

func main() {
	log.Println(uint64(3) * uint64(65) / uint64(5))
	a4 := uint64(3) * uint64(32) / uint64(5)
	a5 := uint64(3) * uint64(33) / uint64(5)

	log.Println("1", a4)
	log.Println("1", a5)
	log.Println("1", a4+a5)

}
