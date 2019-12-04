package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var period = flag.Duration("period", 1*time.Second, "sleep period")

func init() {
	flag.Parse()
	fmt.Printf("Sleeping for %v...\n", *period)
}

func main() {
	for idx, args := range os.Args {
		fmt.Println("参数"+strconv.Itoa(idx)+":", args)
	}
	time.Sleep(*period)
	fmt.Println("end")
}
