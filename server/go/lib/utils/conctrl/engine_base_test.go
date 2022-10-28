package conctrl

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	engine := NewEngine(10)

	tasks := make([]*Task, 100)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = taskgen(strconv.Itoa(i), engine)
	}
	engine.Run(tasks...)
}

func taskgen(id string, engine *BaseEngine) *Task {
	return &Task{Do: func(ctx context.Context) {
		log.Println("task", id)
		n := rand.Intn(10)
		//log.Println("rand", n)
		if n < 5 {
			for i := 0; i < n; i++ {
				engine.AddTask(taskgen(id+"_"+strconv.Itoa(i), engine))
			}
		}
	}}
}

func TestEngineOneTask(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	engine := NewEngine(10)
	ch := make(chan string)
	go func() {
		for {
			id := <-ch
			log.Println("rand", id)
			for i := 0; i < 3; i++ {
				engine.AddTask(taskgen2(id+"_"+strconv.Itoa(i), ch))
			}
		}
	}()
	engine.Run(&Task{
		Do: func(ctx context.Context) {
			ch <- "1"
		},
	})
}

func taskgen2(id string, ch chan string) *Task {
	return &Task{
		Do: func(ctx context.Context) {
			log.Println("task", id)
			n := rand.Intn(10)
			//log.Println("rand", n)
			if n < 3 {
				ch <- id
			}
		},
	}
}
