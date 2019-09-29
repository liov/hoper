package main

import (
	"fmt"
	"time"
)

func main() {
	defer timeCost(time.Now())
	fmt.Println("start program")
	time.Sleep(5 * time.Second)
	fmt.Println("finish program")
}

func timeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println(terminal)
}
