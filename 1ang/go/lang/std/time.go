package main

import (
	"log"
	"time"
)

func main() {
	t := time.Time{}
	log.Println(t.Local())
	log.Println(t.Location())
}
