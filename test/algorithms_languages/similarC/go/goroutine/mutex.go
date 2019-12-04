package main

import (
	"fmt"
	"sync"
	"time"
)

type Book struct {
	BookName string
	L        *sync.Mutex
}

func (bk *Book) SetName(wg *sync.WaitGroup, name string) {
	defer func() {
		fmt.Println("Unlock set name:", name)
		bk.L.Unlock()
		wg.Done()
	}()

	bk.L.Lock()
	fmt.Println("Lock set name:", name)
	time.Sleep(1 * time.Second)
	bk.BookName = name
}

func main() {
	c := sync.NewCond(&sync.Mutex{})    //1
	queue := make([]interface{}, 0, 10) //2

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()        //8
		queue = queue[1:] //9
		fmt.Println("Removed from queue")
		c.L.Unlock() //10
		c.Signal()   //11
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()            //3
		for len(queue) == 2 { //4
			c.Wait() //5
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second) //6
		c.L.Unlock()                        //7
	}

	wg := sync.WaitGroup{}

	var mutex sync.Mutex
	fmt.Println("Locking  (G0)")
	mutex.Lock()
	fmt.Println("locked (G0)")
	wg.Add(3)

	for i := 1; i < 4; i++ {
		go func(i int) {
			fmt.Printf("Locking (G%d)\n", i)
			mutex.Lock()
			fmt.Printf("locked (G%d)\n", i)

			time.Sleep(time.Second * 2)
			mutex.Unlock()
			fmt.Printf("unlocked (G%d)\n", i)
			wg.Done()
		}(i)
	}

	time.Sleep(time.Second * 5)
	fmt.Println("ready unlock (G0)")
	mutex.Unlock()
	fmt.Println("unlocked (G0)")
	wg.Wait()
}

func mutexStruct() {

	bk := Book{}
	bk.L = new(sync.Mutex)
	wg := &sync.WaitGroup{}
	books := []string{"《三国演义》", "《道德经》", "《西游记》"}
	for _, book := range books {
		wg.Add(1)
		go bk.SetName(wg, book)
	}

	wg.Wait()
}
