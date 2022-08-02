package conctrl

import (
	"context"
	"log"
	"testing"
)

func TestEngine(t *testing.T) {
	engine := NewEngine(30)
	taskgen := func(id int) func(ctx context.Context) {
		return func(ctx context.Context) {
			log.Println("task", id)
		}
	}
	tasks := make([]Task, 100)
	for i := 0; i < 100; i++ {
		tasks[i] = taskgen(i)
	}
	engine.Run(tasks...)
}
