package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/structure/heap"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/structure/list"
	synci "github.com/actliboy/hoper/server/go/lib/utils/sync"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type BaseEngine[T, W any] struct {
	limitWorkerCount, currentWorkerCount uint64
	limitWaitTaskCount                   uint
	workerChan                           chan *Worker[T, W]
	taskChan                             chan *BaseTask[T]
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
	EngineStatistics
}

type EngineStatistics struct {
	WorkStatistics
}

func NewBaseEngine[T, W any](workerCount uint) *BaseEngine[T, W] {
	return NewBaseEngineWithContext[T, W](workerCount, context.Background())
}

func NewBaseEngineWithContext[T, W any](workerCount uint, ctx context.Context) *BaseEngine[T, W] {
	ctx, cancel := context.WithCancel(ctx)
	return &BaseEngine[T, W]{
		limitWorkerCount:   uint64(workerCount),
		limitWaitTaskCount: workerCount * 10,
		ctx:                ctx,
		cancel:             cancel,
		workerChan:         make(chan *Worker[T, W]),
		taskChan:           make(chan *BaseTask[T]),
	}
}

func (e *BaseEngine[T, W]) Context() context.Context {
	return e.ctx
}

func (e *BaseEngine[T, W]) Cancel() {
	log.Println("任务取消")
	e.cancel()
	synci.WaitGroupStopWait(&e.wg)

}

func (e *BaseEngine[T, W]) Run(tasks ...*BaseTask[T]) {
	e.addWorker()

	go func() {
		timer := time.NewTimer(time.Second * 1)
		workerList := list.NewSimpleList[*Worker[T, W]]()
		taskList := heap.Heap[*BaseTask[T]]{}
		var emptyTimes, stopTimes uint
		var readyWorkerCh chan *BaseTask[T]
		var readyTask *BaseTask[T]
	loop:
		for {
			if workerList.Size > 0 && len(taskList) > 0 {
				if readyWorkerCh == nil {
					readyWorkerCh = workerList.First().taskCh
				}
				if readyTask == nil {
					readyTask = taskList.First()
				}
			}
			if len(taskList) > int(e.limitWaitTaskCount) {
				select {
				case readyWorker := <-e.workerChan:
					workerList.Push(readyWorker)
				case readyWorkerCh <- readyTask:
					workerList.Pop()
					taskList.Pop()
					readyWorkerCh = nil
					readyTask = nil
				case <-timer.C:
					//检测任务是否卡住
					stopTimes++
					if stopTimes == 10 {
						e.limitWaitTaskCount += uint(e.limitWorkerCount)
						stopTimes = 0
					}
					timer.Reset(time.Second * 1)
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
					readyWorkerCh = nil
					readyTask = nil
				case <-timer.C:
					//检测任务是否已空
					if workerList.Size == uint(e.currentWorkerCount) && len(taskList) == 0 {
						emptyTimes++
						if emptyTimes > 2 {
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

func (e *BaseEngine[T, W]) newWorker(readyTask *BaseTask[T]) {
	e.currentWorkerCount++
	//id := c.currentWorkerCount
	taskChan := make(chan *BaseTask[T])
	worker := &Worker[T, W]{Id: uint(e.currentWorkerCount), taskCh: taskChan}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				log.Println(string(debug.Stack()))
				e.wg.Done()
			}
			atomic.AddUint64(&e.currentWorkerCount, ^uint64(0))
		}()
		if readyTask != nil && readyTask.BaseTaskFunc != nil {
			readyTask.BaseTaskFunc(e.ctx)
			e.wg.Done()
		}
		for {
			select {
			case e.workerChan <- worker:
				task := <-taskChan
				if task != nil && task.BaseTaskFunc != nil {
					task.BaseTaskFunc(e.ctx)
					atomic.AddUint64(&e.taskDoneCount, 1)
				}
				e.wg.Done()
			case <-e.ctx.Done():
				return
			}
		}
	}()
}

func (e *BaseEngine[T, W]) addWorker() {
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

func (e *BaseEngine[T, W]) AddTask(task *BaseTask[T]) {
	if task == nil || task.BaseTaskFunc == nil {
		return
	}
	atomic.AddUint64(&e.taskTotalCount, 1)
	e.wg.Add(1)
	e.taskChan <- task
}

func (e *BaseEngine[T, W]) AddTasks(tasks ...*BaseTask[T]) {
	e.wg.Add(len(tasks))
	for _, task := range tasks {
		e.taskChan <- task
	}
}

func (e *BaseEngine[T, W]) AddWorker(num int) {
	atomic.AddUint64(&e.limitWorkerCount, uint64(num))
	e.addWorker()
}
