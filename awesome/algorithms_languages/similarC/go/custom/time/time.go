package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().Year())
	fmt.Println(int(time.Now().Month()))
	fmt.Println(time.Now().Day())
	fmt.Println(time.Now().Format("2006-01-02T15:04:05Z+08:00"))
}
