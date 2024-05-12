package main

import (
	"fmt"
	"time"
)

func main() {
	t1, _ := time.ParseInLocation(time.DateTime, "0000-00-00 00:00:00", time.Local)
	t2, _ := time.Parse(time.DateTime, "0000-00-00 00:00:00")
	fmt.Println(t1, t2)
	fmt.Println(t1.IsZero(), t2.IsZero())
}
