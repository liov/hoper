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
	go func() {
		time.AfterFunc(time.Second*2, func() {
			c.taskChan <- func(ctx context.Context) {
				log.Println("额外任务", c.wg)
			}
		})
	}()
	c.addWorker()

	go func() {
		workerList := list.NewSimpleList[*Worker]()
		taskList := list.NewSimpleList[Task]()

	loop:
		for {

			var readyWorkerCh chan Task
			var readyTask Task
			if workerList.Size > 0 && taskList.Size > 0 {
				readyWorker := workerList.Pop()
				log.Println("ready worker :", readyWorker.id)
				readyWorkerCh = readyWorker.ch
				readyTask = taskList.Pop()
			}
			select {
			case readyTask = <-c.taskChan:
				c.wg.Add(1)
				taskList.Push(readyTask)
			case readyWorker := <-c.workerChan:
				workerList.Push(readyWorker)
			case readyWorkerCh <- readyTask:
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
	//id := c.currentWorkerCount
	taskChan := make(chan Task)
	worker := &Worker{c.currentWorkerCount, taskChan}
	go func() {
		for {
			select {
			case c.workerChan <- worker:
				log.Println("worker id :", worker.id, "worker count :", c.currentWorkerCount)
				task := <-taskChan
				log.Println("有任务做了")
				if task != nil {
					task(c.ctx)
					//log.Println("task is done,worker id :", id) 删掉就会报错
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
	c.taskChan <- task
}

func (c *Engine) AddWorker(num int) {
	atomic.AddInt64(&c.limitWorkerCount, int64(num))
	c.addWorker()
}
