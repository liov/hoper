package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/structure/list"
	"log"
	"sync"
	"sync/atomic"
)

type Task func(context.Context)

type Engine struct {
	limitWorkerCount, currentWorkerCount int64
	workerChan                           chan chan Task
	taskChan                             chan Task
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
}

func NewEngine(workerCount int) *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine{
		limitWorkerCount: int64(workerCount),
		ctx:              ctx,
		cancel:           cancel,
		workerChan:       make(chan chan Task),
		taskChan:         make(chan Task),
	}
}

func (c *Engine) Run(tasks ...Task) {

	c.addWorker()

	go func() {
		workerList := list.NewSimpleList[chan Task]()
		requestList := list.NewSimpleList[Task]()
	loop:
		for {
			var readyWorker chan Task
			var readyTask Task
			if workerList.Size > 0 && requestList.Size > 0 {
				readyWorker = workerList.Pop()
				readyTask = requestList.Pop()
			}
			select {
			case readyTask = <-c.taskChan:
				requestList.Push(readyTask)
			case readyWorker = <-c.workerChan:
				workerList.Push(readyWorker)
			case readyWorker <- readyTask:
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

func (c *Engine) newWorker() {
	c.currentWorkerCount++
	id := c.currentWorkerCount
	taskChan := make(chan Task)
	go func() {
		for {
			select {
			case c.workerChan <- taskChan:
				task := <-taskChan
				if task != nil {
					task(c.ctx)
					log.Println("task is done,worker id :", id)
				}
				c.wg.Done()
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *Engine) addWorker() {
	if c.currentWorkerCount == 0 {
		c.newWorker()
	}
	go func() {
		for {
			select {
			case readyTask := <-c.taskChan:
				if c.currentWorkerCount < c.limitWorkerCount {
					c.newWorker()
				}
				c.taskChan <- readyTask
				if c.currentWorkerCount == c.limitWorkerCount {
					return
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *Engine) AddTask(task Task) {
	c.wg.Add(1)
	c.taskChan <- task
}

func (c *Engine) AddWorker(num int) {
	atomic.AddInt64(&c.limitWorkerCount, int64(num))
	c.addWorker()
}

type Worker struct {
	idleChan chan struct{}
}
