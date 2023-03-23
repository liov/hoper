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
	engine := NewEngine[int, int, int](10)

	tasks := make([]*Task[int, int], 100)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = taskGen(strconv.Itoa(i), engine)
	}
	engine.Run(tasks...)
}

func taskGen(id string, engine *Engine[int, int, int]) *Task[int, int] {
	return &Task[int, int]{
		TaskFunc: func(ctx context.Context) ([]*Task[int, int], error) {
			log.Println("task", id)
			n := rand.Intn(10)
			//log.Println("rand", n)
			if n < 3 {
				for i := 0; i < n; i++ {
					engine.BaseEngine.AddTask(engine.BaseTask(taskGen(id+"_"+strconv.Itoa(i), engine)))
				}
			}
			if n == 3 {
				panic(n)
			}
			return nil, nil
		}}
}
