package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/slices"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/structure/list"
	synci "github.com/actliboy/hoper/server/go/lib/utils/sync"
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

// TODO
type TaskMeta struct {
	Id   uint
	Kind Kind
}

type TaskStatistics struct {
	timeCost time.Duration
}

type Task struct {
	TaskMeta
	Do TaskFun
}

type ErrHandle func(context.Context, error)

// TaskWithErrHandle Deprecated
// 原本设计框架参与error处理，但是error处理仍然需要传参指定，不如就在task内部自己处理掉
type TaskWithErrHandle struct {
	TaskMeta
	Do        func(context.Context) error
	ErrHandle ErrHandle
}

type Worker struct {
	Id     uint
	Kind   Kind
	taskCh chan *Task
}

type Engine struct {
	limitWorkerCount, currentWorkerCount uint64
	workerChan                           chan *Worker
	taskChan                             chan *Task
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
	kindHandler                          []KindHandler
	averageTimeCost                      time.Duration
	taskDoneCount, taskTotalCount        uint64
}

type KindHandler struct {
	Skip bool
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

func (e *Engine) Context() context.Context {
	return e.ctx
}

func (e *Engine) Cancel() {
	log.Println("任务取消")
	e.cancel()
	synci.WaitGroupStopWait(&e.wg)

}

func (e *Engine) SkipKind(kinds ...Kind) *Engine {
	length := slices.Max(kinds) + 1
	if e.kindHandler == nil {
		e.kindHandler = make([]KindHandler, length)
	}
	if int(length) > len(e.kindHandler) {
		e.kindHandler = append(e.kindHandler, make([]KindHandler, int(length)-len(e.kindHandler))...)
	}
	for _, kind := range kinds {
		e.kindHandler[kind].Skip = true
	}
	return e
}

func (e *Engine) Run(tasks ...*Task) {
	e.addWorker()

	go func() {
		timer := time.NewTimer(time.Second * 1)
		workerList := list.NewSimpleList[*Worker]()
		taskList := list.NewSimpleList[*Task]()
		var emptyTimes int
	loop:
		for {
			var readyWorkerCh chan *Task
			var readyTask *Task
			if workerList.Size > 0 && taskList.Size > 0 {
				readyWorkerCh = workerList.First().taskCh
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
				case <-timer.C:
					//检测任务是否已空
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
					timer.Stop()
					break loop
				}
			}

		}
	}()

	e.taskTotalCount = uint64(len(tasks))
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
	worker := &Worker{Id: uint(e.currentWorkerCount), taskCh: taskChan}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				log.Println(string(debug.Stack()))
				e.wg.Done()
			}
			atomic.AddUint64(&e.currentWorkerCount, ^uint64(0))
		}()
		if readyTask != nil && readyTask.Do != nil {
			readyTask.Do(e.ctx)
			e.wg.Done()
		}
		for {
			select {
			case e.workerChan <- worker:
				task := <-taskChan
				if task != nil && task.Do != nil {
					task.Do(e.ctx)
					atomic.AddUint64(&e.taskDoneCount, 1)
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
	if task == nil {
		return
	}
	if e.kindHandler != nil && int(task.Kind) < len(e.kindHandler) && e.kindHandler[task.Kind].Skip {
		return
	}
	atomic.AddUint64(&e.taskTotalCount, 1)
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
