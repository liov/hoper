package main

import (
	"fmt"
	"time"

	"github.com/liov/hoper/server/go/lib/utils/struct/cache/gcache"
)

func main() {
	gc := gcache.New(10).
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
