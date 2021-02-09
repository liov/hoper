package main

import (
	"fmt"

	"time"

	"github.com/liov/hoper/go/v2/utils/dao/cache"
)

func main() {
	gc := cache.New(10).
		LFU().
		Build()

	gc.SetWithExpire("key", "ok", time.Second*3)

	v, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("value:", v)

	fmt.Println("waiting 3s for value to expire:")
	time.Sleep(time.Second * 3)

	v, err = gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("value:", v)
}
