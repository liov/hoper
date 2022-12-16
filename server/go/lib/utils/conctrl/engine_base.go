package conctrl

import (
	"context"
	synci "github.com/liov/hoper/server/go/lib/utils/sync"
	"github.com/liov/hoper/server/go/lib_v2/utils/structure/heap"
	"github.com/liov/hoper/server/go/lib_v2/utils/structure/list"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type BaseEngine struct {
	limitWorkerCount, currentWorkerCount uint64
	limitWaitTaskCount                   uint
	workerChan                           chan *Worker
	taskChan                             chan *BaseTask
	ctx                                  context.Context
	cancel                               context.CancelFunc
	wg                                   sync.WaitGroup
	done                                 sync.Map
	kindHandler                          []*KindHandler
	EngineStatistics
}

type EngineStatistics struct {
	WorkStatistics
}

func NewBaseEngine(workerCount uint) *BaseEngine {
	return NewBaseEngineWithContext(workerCount, context.Background())
}

func NewBaseEngineWithContext(workerCount uint, ctx context.Context) *BaseEngine {
	ctx, cancel := context.WithCancel(ctx)
	return &BaseEngine{
		limitWorkerCount:   uint64(workerCount),
		limitWaitTaskCount: workerCount * 10,
		ctx:                ctx,
		cancel:             cancel,
		workerChan:         make(chan *Worker),
		taskChan:           make(chan *BaseTask),
	}
}

func (e *BaseEngine) Context() context.Context {
	return e.ctx
}

func (e *BaseEngine) Cancel() {
	log.Println("任务取消")
	e.cancel()
	synci.WaitGroupStopWait(&e.wg)

}

func (e *BaseEngine) Run(tasks ...*BaseTask) {
	e.addWorker()

	go func() {
		timer := time.NewTimer(time.Second * 5)
		workerList := list.NewSimpleList[*Worker]()
		taskList := heap.Heap[*BaseTask]{}
		var emptyTimes uint
		var readyWorkerCh chan *BaseTask
		var readyTask *BaseTask
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

func (e *BaseEngine) newWorker(readyTask *BaseTask) {
	e.currentWorkerCount++
	//id := c.currentWorkerCount
	taskChan := make(chan *BaseTask)
	worker := &Worker{Id: uint(e.currentWorkerCount), taskCh: taskChan}
	go func() {
		defer func() {
			if r := recover(); r != nil {
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
	if task == nil || task.BaseTaskFunc == nil {
		return
	}
	atomic.AddUint64(&e.taskTotalCount, 1)
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
