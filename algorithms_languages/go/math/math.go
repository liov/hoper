package main

import (
	"log"
	"math/big"
)

func main() {
	f1 := big.NewFloat(5.23)
	f2 := big.NewFloat(0.1000000000000000055511)
	log.Println(f2)
	f3 := &big.Float{}
	f3 = f3.Add(f1, f2)
	log.Println(f3)
}
