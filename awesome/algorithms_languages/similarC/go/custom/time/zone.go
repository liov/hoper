package main

import (
	"log"
	"time"
)

func main() {
	now := time.Now()
	now = now.UTC().AddDate(0, 0, 1)
	log.Println(now)
}
