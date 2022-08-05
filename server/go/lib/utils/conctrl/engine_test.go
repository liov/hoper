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

	tasks := make([]Task, 10)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = taskgen(strconv.Itoa(i), engine)
	}
	engine.Run(tasks...)
}

func taskgen(id string, engine *Engine) func(ctx context.Context) {
	return func(ctx context.Context) {
		log.Println("task", id)
		n := rand.Intn(10)
		if n < 3 {
			for i := 0; i < n; i++ {
				engine.AddTask(taskgen(id+"_"+strconv.Itoa(i), engine))
			}
		}
	}
}
