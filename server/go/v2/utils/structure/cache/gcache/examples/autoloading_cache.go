package main

import (
	"fmt"

	"github.com/liov/hoper/go/v2/utils/dao/cache"
)

func main() {
	gc := cache.New(10).
		LFU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
		return fmt.Sprintf("%v-value", key), nil
	}).
		Build()

	v, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
