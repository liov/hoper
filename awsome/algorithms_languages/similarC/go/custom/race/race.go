package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	ac           = make(chan string, 10)
	openclose    = make(chan bool)
	continueread = true
	arr          = []string{}
)

type Test struct {
	mu sync.Mutex
}

func (t *Test) AutoCommit() {
	tc := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-tc.C:
			log.Println("AutoCommit进来了")
			t.mu.Lock()
			// aa := make([]string, len(arr))
			// copy(aa, arr)
			openclose <- false
			arr = []string{}
			t.mu.Unlock()
			log.Println("AutoCommit出去了")
		}
	}
}

func (t *Test) DataOp(item string) {
	t.mu.Lock()
	fmt.Println(item)
	t.mu.Unlock()
}

func main() {
	test := new(Test)

	go func() {
		for {
			continueread = <-openclose
			log.Println(continueread)
			fmt.Println("read paused")
		}
	}()

	go func() {
		tc := time.NewTicker(time.Second * 5)
		for range tc.C {
			log.Println("进来了")
			if !continueread {
				continueread = true
				fmt.Println("read resumed")
			}
		}
	}()

	go func() {
		index := 0
		for {
			if continueread {
				ac <- fmt.Sprintf("%d", index)
				index++
			}
		}
	}()

	go test.AutoCommit()

	go func() {
		for v := range ac {
			arr = append(arr, v)
			test.DataOp(v)
		}
	}()

	c := make(chan struct{})
	<-c
}
