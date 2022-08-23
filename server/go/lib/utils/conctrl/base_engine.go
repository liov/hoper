package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/structure/list"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type BaseTask struct {
	Id uint
	Do TaskFun
}

type BaseWorker struct {
	Id uint
	ch chan *BaseTask
}

type BaseEngine struct {
	limitWorkerCount, currentWorkerCount uint64
	workerChan                           chan *BaseWorker
	taskChan                             chan *BaseTask
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
}

func NewBaseEngine(workerCount uint) *BaseEngine {
	return NewBaseEngineWithContext(workerCount, context.Background())
}

func NewBaseEngineWithContext(workerCount uint, ctx context.Context) *BaseEngine {
	ctx, cancel := context.WithCancel(ctx)
	return &BaseEngine{
		limitWorkerCount: uint64(workerCount),
		ctx:              ctx,
		cancel:           cancel,
		workerChan:       make(chan *BaseWorker),
		taskChan:         make(chan *BaseTask),
	}
}

func (e *BaseEngine) Run(tasks ...*BaseTask) {
	e.addWorker()

	go func() {
		workerList := list.NewSimpleList[*BaseWorker]()
		taskList := list.NewSimpleList[*BaseTask]()
		timer := time.NewTimer(time.Second * 1)
		var emptyTimes int
	loop:
		for {
			var readyWorkerCh chan *BaseTask
			var readyTask *BaseTask
			if workerList.Size > 0 && taskList.Size > 0 {
				readyWorkerCh = workerList.First().ch
				readyTask = taskList.First()
			}
			if taskList.Size > int(e.limitWorkerCount*2) {
				select {
				case readyWorker := <-e.workerChan:
					workerList.Push(readyWorker)
				case readyWorkerCh <- readyTask:
					workerList.Pop()
					taskList.Pop()
				case <-e.ctx.Done():
					timer.Stop()
					break loop
				}
			} else {
				select {
				case readyTask = <-e.taskChan:
					taskList.Push(readyTask)
				case readyWorker := <-e.workerChan:
					workerList.Push(readyWorker)
				case readyWorkerCh <- readyTask:
					workerList.Pop()
					taskList.Pop()
					//检测任务是否已空
				case <-timer.C:
					if workerList.Size == int(e.currentWorkerCount) && taskList.Size == 0 {
						emptyTimes++
						if emptyTimes > 2 {
							log.Println("task is empty")
							log.Println("任务即将结束")
							e.wg.Done()
							timer.Stop()
							break loop
						}
					}
					emptyTimes = 0
					timer.Reset(time.Second * 1)
				case <-e.ctx.Done():
					e.wg.Done()
					timer.Stop()
					break loop
				}
			}
		}
	}()

	e.wg.Add(len(tasks) + 1)
	for _, task := range tasks {
		e.taskChan <- task
	}

	e.wg.Wait()
	log.Println("任务结束")
}

func (e *BaseEngine) newWorker(readyTask *BaseTask) {
	e.currentWorkerCount++
	//id := c.currentWorkerCount
	taskChan := make(chan *BaseTask)
	worker := &BaseWorker{uint(e.currentWorkerCount), taskChan}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				log.Println(string(debug.Stack()))
				e.wg.Done()
			}
			atomic.AddUint64(&e.currentWorkerCount, ^uint64(0))
		}()
		if readyTask != nil {
			readyTask.Do(e.ctx)
			e.wg.Done()
		}
		for {
			select {
			case e.workerChan <- worker:
				task := <-taskChan
				if task != nil {
					task.Do(e.ctx)
				}
				e.wg.Done()
			case <-e.ctx.Done():
				return
			}
		}
	}()
}

func (e *BaseEngine) addWorker() {
	if e.currentWorkerCount == 0 {
		e.newWorker(nil)
	}
	go func() {
		for {
			select {
			case readyTask := <-e.taskChan:
				if e.currentWorkerCount < e.limitWorkerCount {
					e.newWorker(readyTask)
				}
				if e.currentWorkerCount == e.limitWorkerCount {
					log.Println("worker count is full")
					return
				}
			case <-e.ctx.Done():
				return
			}
		}
	}()
}

func (e *BaseEngine) AddTask(task *BaseTask) {
	if task == nil {
		return
	}
	e.wg.Add(1)
	e.taskChan <- task
}

func (e *BaseEngine) AddTasks(tasks ...*BaseTask) {
	e.wg.Add(len(tasks))
	for _, task := range tasks {
		e.taskChan <- task
	}
}

func (e *BaseEngine) AddWorker(num int) {
	atomic.AddUint64(&e.limitWorkerCount, uint64(num))
	e.addWorker()
}

func (e *BaseEngine) Cancel() {
	e.cancel()
}
