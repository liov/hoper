package main

import "log"

func main() {
	var s = make([]int, 10)
	for i := 0; i < len(s); i++ {
		log.Println(s)
		log.Println(i)
	}
}
