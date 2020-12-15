package main

import (
	"fmt"
	"time"
)

func main() {
	defer timeCost(time.Now())
	defer func(start time.Time) { fmt.Println("这种计时有意义吗", time.Since(start)) }(time.Now())

	fmt.Println("start program")
	time.Sleep(5 * time.Second)
	fmt.Println("finish program")
}

func timeCost(start time.Time) {
	fmt.Println(time.Since(start))
}
