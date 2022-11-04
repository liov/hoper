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

type BaseEngine[KEY comparable, T, W any] struct {
	limitWorkerCount, currentWorkerCount uint64
	limitWaitTaskCount                   uint
	workerChan                           chan *Worker[KEY, T, W]
	taskChan                             chan *BaseTask[KEY, T]
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
	fixedWorker                          []chan *BaseTask[KEY, T]
	EngineStatistics
}

type EngineStatistics struct {
	WorkStatistics
}

func NewBaseEngine[KEY comparable, T, W any](workerCount uint) *BaseEngine[KEY, T, W] {
	return NewBaseEngineWithContext[KEY, T, W](workerCount, context.Background())
}

func NewBaseEngineWithContext[KEY comparable, T, W any](workerCount uint, ctx context.Context) *BaseEngine[KEY, T, W] {
	ctx, cancel := context.WithCancel(ctx)
	return &BaseEngine[KEY, T, W]{
		limitWorkerCount:   uint64(workerCount),
		limitWaitTaskCount: workerCount * 10,
		ctx:                ctx,
		cancel:             cancel,
		workerChan:         make(chan *Worker[KEY, T, W]),
		taskChan:           make(chan *BaseTask[KEY, T]),
	}
}

func (e *BaseEngine[KEY, T, W]) Context() context.Context {
	return e.ctx
}

func (e *BaseEngine[KEY, T, W]) Cancel() {
	log.Println("任务取消")
	e.cancel()
	synci.WaitGroupStopWait(&e.wg)

}

func (e *BaseEngine[KEY, T, W]) Run(tasks ...*BaseTask[KEY, T]) {
	e.addWorker()

	go func() {
		timer := time.NewTimer(time.Second * 5)
		workerList := list.NewSimpleList[*Worker[KEY, T, W]]()
		taskList := heap.Heap[*BaseTask[KEY, T]]{}
		var emptyTimes uint
		var readyWorkerCh chan *BaseTask[KEY, T]
		var readyTask *BaseTask[KEY, T]
	loop:
		for {
			if workerList.Size > 0 && len(taskList) > 0 {
				if readyWorkerCh == nil {
					readyWorkerCh = workerList.Pop().taskCh
				}
				if readyTask == nil {
					readyTask = taskList.Pop()
				}
			}

			if len(taskList) > int(e.limitWaitTaskCount) {
				select {
				case readyWorker := <-e.workerChan:
					workerList.Push(readyWorker)
				case readyWorkerCh <- readyTask:
					readyWorkerCh = nil
					readyTask = nil
				case <-e.ctx.Done():
					timer.Stop()
					break loop
				}
			} else {
				select {
				case readyTaskTmp := <-e.taskChan:
					taskList.Push(readyTaskTmp)
				case readyWorker := <-e.workerChan:
					workerList.Push(readyWorker)
				case readyWorkerCh <- readyTask:
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

func (e *BaseEngine[KEY, T, W]) newWorker(readyTask *BaseTask[KEY, T]) {
	e.currentWorkerCount++
	//id := c.currentWorkerCount
	taskChan := make(chan *BaseTask[KEY, T])
	worker := &Worker[KEY, T, W]{Id: uint(e.currentWorkerCount), taskCh: taskChan}
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
				readyTask = <-taskChan
				if readyTask != nil && readyTask.BaseTaskFunc != nil {
					readyTask.BaseTaskFunc(e.ctx)
					atomic.AddUint64(&e.taskDoneCount, 1)
				}
				e.wg.Done()
			case <-e.ctx.Done():
				return
			}
		}
	}()
}

func (e *BaseEngine[KEY, T, W]) addWorker() {
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

func (e *BaseEngine[KEY, T, W]) AddTask(task *BaseTask[KEY, T]) {
	if task == nil || task.BaseTaskFunc == nil {
		return
	}
	atomic.AddUint64(&e.taskTotalCount, 1)
	e.wg.Add(1)
	e.taskChan <- task
}

func (e *BaseEngine[KEY, T, W]) AddTasks(tasks ...*BaseTask[KEY, T]) {
	e.wg.Add(len(tasks))
	for _, task := range tasks {
		e.taskChan <- task
	}
}

func (e *BaseEngine[KEY, T, W]) AddWorker(num int) {
	atomic.AddUint64(&e.limitWorkerCount, uint64(num))
	e.addWorker()
}

func (e *BaseEngine[KEY, T, W]) NewFixedWorker() int {
	ch := make(chan *BaseTask[KEY, T])
	e.fixedWorker = append(e.fixedWorker, ch)
	go func() {
		for task := range ch {
			task.BaseTaskFunc(e.ctx)
			e.wg.Done()
		}
	}()
	return len(e.fixedWorker) - 1
}

func (e *BaseEngine[KEY, T, W]) AddFixedTask(workerId int, task *BaseTask[KEY, T]) {
	if workerId > len(e.fixedWorker)-1 {
		return
	}
	ch := e.fixedWorker[workerId]
	e.wg.Add(1)
	go func() {
		ch <- task
	}()
}
