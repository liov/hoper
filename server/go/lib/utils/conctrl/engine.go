package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/slices"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/structure/list"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type Kind uint8

const (
	KindNormal = iota
)

type TaskFun func(context.Context)

type Task struct {
	id   uint
	Kind Kind
	Do   TaskFun
}

type ErrHandle func(context.Context, error)

// TaskWithErrHandle Deprecated
// 原本设计框架参与error处理，但是error处理仍然需要传参指定，不如就在task内部自己处理掉
type TaskWithErrHandle struct {
	id        uint
	Kind      Kind
	Do        func(context.Context) error
	ErrHandle ErrHandle
}

type Worker struct {
	id uint
	ch chan *Task
}

type Engine struct {
	limitWorkerCount, currentWorkerCount uint64
	workerChan                           chan *Worker
	taskChan                             chan *Task
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
	excludeKinds                         []bool
}

func NewEngine(workerCount uint) *Engine {
	return NewEngineWithContext(workerCount, context.Background())
}

func NewEngineWithContext(workerCount uint, ctx context.Context) *Engine {
	ctx, cancel := context.WithCancel(ctx)
	return &Engine{
		limitWorkerCount: uint64(workerCount),
		ctx:              ctx,
		cancel:           cancel,
		workerChan:       make(chan *Worker),
		taskChan:         make(chan *Task),
	}
}

func (e *Engine) SkipKind(kinds ...Kind) *Engine {
	length := slices.Max(kinds) + 1
	if e.excludeKinds == nil {
		e.excludeKinds = make([]bool, length)
	}
	if int(length) > len(e.excludeKinds) {
		e.excludeKinds = append(e.excludeKinds, make([]bool, int(length)-len(e.excludeKinds))...)
	}
	for _, kind := range kinds {
		e.excludeKinds[kind] = true
	}
	return e
}

func (e *Engine) Run(tasks ...*Task) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println(string(debug.Stack()))
		}
	}()
	e.addWorker()

	go func() {
		workerList := list.NewSimpleList[*Worker]()
		taskList := list.NewSimpleList[*Task]()
		timer := time.NewTimer(time.Second * 1)
		var emptyTimes int
	loop:
		for {
			var readyWorkerCh chan *Task
			var readyTask *Task
			if workerList.Size > 0 && taskList.Size > 0 {
				readyWorkerCh = workerList.First().ch
				readyTask = taskList.First()
			}
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
					if emptyTimes > 5 {
						log.Println("task is empty")
						log.Println("任务即将结束")
						e.wg.Done()
						timer.Stop()
						break loop
					}
				}
				timer.Reset(time.Second * 1)
			case <-e.ctx.Done():
				timer.Stop()
				break loop
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

func (e *Engine) newWorker(readyTask *Task) {
	e.currentWorkerCount++
	//id := c.currentWorkerCount
	taskChan := make(chan *Task)
	worker := &Worker{uint(e.currentWorkerCount), taskChan}
	go func() {
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

func (e *Engine) addWorker() {
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

func (e *Engine) AddTask(task *Task) {
	if e.excludeKinds != nil && int(task.Kind) < len(e.excludeKinds) && e.excludeKinds[task.Kind] {
		return
	}
	e.wg.Add(1)
	e.taskChan <- task
}

func (e *Engine) AddTasks(tasks ...*Task) {
	e.wg.Add(len(tasks))
	for _, task := range tasks {
		e.taskChan <- task
	}
}

func (e *Engine) AddWorker(num int) {
	atomic.AddUint64(&e.limitWorkerCount, uint64(num))
	e.addWorker()
}