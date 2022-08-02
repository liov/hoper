package conctrl

import (
	"context"
	"log"
	"sync"
)

type Task[T any] func(context.Context) (T, error)

type TaskB[REQ, RES any] func(context.Context, REQ) (RES, error)

type TaskCo[REQ, RES any] interface {
	Task[RES] | TaskB[REQ, RES]
}

type Engine[T any] struct {
	limitWorkerCount, currentWorkerCount int
	workerChan                           chan struct{}
	taskChan                             chan Task[T]
	//taskList   list.SimpleList[Task]
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewEngine[T any](workerCount int) *Engine[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine[T]{
		limitWorkerCount: workerCount,
		workerChan:       make(chan struct{}, workerCount),
		ctx:              ctx,
		cancel:           cancel,
		taskChan:         make(chan Task[T]),
	}
}

func (c *Engine[T]) Run(tasks ...Task[T]) {

	if c.currentWorkerCount == 0 {
		c.NewWorker()
	}

	go func() {
	loop:
		for {
			select {
			case readyTask := <-c.taskChan:
				if c.currentWorkerCount < c.limitWorkerCount {
					c.NewWorker()
				}
				c.taskChan <- readyTask
				if c.currentWorkerCount == c.limitWorkerCount {
					break loop
				}
			case <-c.ctx.Done():
				break loop
			}
		}
	}()
	c.wg.Add(len(tasks))
	for _, task := range tasks {
		c.taskChan <- task
	}

	c.wg.Wait()
}

func (c *Engine[T]) NewWorker() {
	select {
	case c.workerChan <- struct{}{}:
		c.currentWorkerCount++
		id := c.currentWorkerCount
		go func() {
		loop:
			for {
				select {
				case task := <-c.taskChan:
					if task != nil {
						task(c.ctx)
						log.Println("task is done,worker id :", id)
					}
					c.wg.Done()
				case <-c.ctx.Done():
					break loop
				}
			}
		}()
	default:
		log.Println("worker is full")
	}
}

func (c *Engine[T]) AddTask(task Task[T]) {
	c.wg.Add(1)
	c.taskChan <- task
}

type Worker[REQ, RES any, T TaskCo[REQ, RES]] struct {
	isReady bool
}
