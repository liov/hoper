package conctrl

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestBaseEngine(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	engine := NewBaseEngine[int, int, int](10)

	tasks := make([]*BaseTask[int, int], 100)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = baseTaskGen(strconv.Itoa(i), engine)
	}
	engine.Run(tasks...)
}

func baseTaskGen(id string, engine *BaseEngine[int, int, int]) *BaseTask[int, int] {
	return &BaseTask[int, int]{BaseTaskFunc: func(ctx context.Context) {
		log.Println("task", id)
		n := rand.Intn(10)
		//log.Println("rand", n)
		if n < 3 {
			for i := 0; i < n; i++ {
				engine.AddTask(baseTaskGen(id+"_"+strconv.Itoa(i), engine))
			}
		}
		if n == 3 {
			panic(n)
		}
	}}
}

func TestBaseEngineOneTask(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	engine := NewBaseEngine[int, int, int](10)
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
	engine.Run(&BaseTask[int, int]{
		BaseTaskFunc: func(ctx context.Context) {
			ch <- "1"
		},
	})
}

func taskgen2(id string, ch chan string) *BaseTask[int, int] {
	return &BaseTask[int, int]{
		BaseTaskFunc: func(ctx context.Context) {
			log.Println("task", id)
			n := rand.Intn(10)
			//log.Println("rand", n)
			if n == 5 {
				panic("5")
			}
			if n < 3 {
				ch <- id
			}
		},
	}
}
