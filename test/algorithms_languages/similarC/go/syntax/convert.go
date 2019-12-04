package main

import "log"

func main() {
	var id uint64
	id = 10
	log.Println(string(byte(id/1_000_000 + 49)))
}
