package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println(bar())
}

func bar() error {
	return errors.Unwrap(fmt.Errorf("%w",errors.New("UnWarp")))
}
