package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/structure/list"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Task func(context.Context)

type Engine struct {
	limitWorkerCount, currentWorkerCount int64
	workerChan                           chan *Worker
	taskChan                             chan Task
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
}

type Worker struct {
	id int64
	ch chan Task
}

func NewEngine(workerCount int) *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	return &Engine{
		limitWorkerCount: int64(workerCount),
		ctx:              ctx,
		cancel:           cancel,
		workerChan:       make(chan *Worker),
		taskChan:         make(chan Task),
	}
}

func (c *Engine) Run(tasks ...Task) {

	c.addWorker()

	go func() {
		workerList := list.NewSimpleList[*Worker]()
		taskList := list.NewSimpleList[Task]()
		timer := time.NewTimer(time.Second * 1)
		var emptyTimes int
	loop:
		for {
			var readyWorkerCh chan Task
			var readyTask Task
			if workerList.Size > 0 && taskList.Size > 0 {
				readyWorkerCh = workerList.First().ch
				readyTask = taskList.First()
			}
			select {
			case readyTask = <-c.taskChan:
				taskList.Push(readyTask)
			case readyWorker := <-c.workerChan:
				workerList.Push(readyWorker)
			case readyWorkerCh <- readyTask:
				workerList.Pop()
				taskList.Pop()
				//检测任务是否已空
			case <-timer.C:
				if workerList.Size == int(c.currentWorkerCount) && taskList.Size == 0 {
					emptyTimes++
					log.Println("任务即将结束")
					if emptyTimes > 5 {
						log.Println("task is empty")
						c.wg.Done()
					}
				}
				timer.Reset(time.Second * 1)
			case <-c.ctx.Done():
				timer.Stop()
				break loop
			}
		}
	}()

	c.wg.Add(len(tasks) + 1)
	for _, task := range tasks {
		c.taskChan <- task
	}

	c.wg.Wait()
}

func (c *Engine) newWorker(readyTask Task) {
	c.currentWorkerCount++
	//id := c.currentWorkerCount
	taskChan := make(chan Task)
	worker := &Worker{c.currentWorkerCount, taskChan}
	go func() {
		if readyTask != nil {
			readyTask(c.ctx)
			c.wg.Done()
		}
		for {
			select {
			case c.workerChan <- worker:
				task := <-taskChan
				if task != nil {
					task(c.ctx)
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
		c.newWorker(nil)
	}
	go func() {
		for {
			select {
			case readyTask := <-c.taskChan:
				if c.currentWorkerCount < c.limitWorkerCount {
					c.newWorker(readyTask)
				}
				if c.currentWorkerCount == c.limitWorkerCount {
					log.Println("worker count is full")
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

func (c *Engine) AddTasks(tasks ...Task) {
	c.wg.Add(len(tasks))
	for _, task := range tasks {
		c.taskChan <- task
	}
}

func (c *Engine) AddWorker(num int) {
	atomic.AddInt64(&c.limitWorkerCount, int64(num))
	c.addWorker()
}
