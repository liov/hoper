package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func main() {
	name(50)
}

type Queue interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func name[T Queue](x T) {
	constraints.Float()
	fmt.Println(x)
}
