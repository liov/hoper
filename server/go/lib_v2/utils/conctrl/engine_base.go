package conctrl

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	rate2 "github.com/liov/hoper/server/go/lib/utils/conctrl/rate"
	"github.com/liov/hoper/server/go/lib/utils/gen"
	synci "github.com/liov/hoper/server/go/lib/utils/sync"
	"github.com/liov/hoper/server/go/lib/v2/utils/structure/heap"
	"github.com/liov/hoper/server/go/lib/v2/utils/structure/list"
	"golang.org/x/time/rate"
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
	workers                              []*Worker[KEY, T, W]
	workerList                           list.SimpleList[*Worker[KEY, T, W]]
	taskChan                             chan *BaseTask[KEY, T]
	taskList                             heap.Heap[*BaseTask[KEY, T]]
	ctx                                  context.Context
	cancel                               context.CancelFunc       // 手动停止执行
	wg                                   sync.WaitGroup           // 控制确保所有任务执行完
	fixedWorker                          []chan *BaseTask[KEY, T] // 固定只执行一种任务的worker,避免并发问题
	speedLimit                           *rate2.SpeedLimiter
	rateLimiter                          *rate.Limiter
	//TODO
	monitor *time.Ticker // 全局检测定时器，任务的卡住检测，worker panic recover都可以用这个检测
	EngineStatistics
	isRunning, isFinished bool
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
		workerList:         list.NewSimpleList[*Worker[KEY, T, W]](),
		taskList:           heap.Heap[*BaseTask[KEY, T]]{},
	}
}

func (e *BaseEngine[KEY, T, W]) Context() context.Context {
	return e.ctx
}

func (e *BaseEngine[KEY, T, W]) SpeedLimited(interval time.Duration) {
	e.speedLimit = rate2.NewSpeedLimiter(interval)
}

func (e *BaseEngine[KEY, T, W]) RandSpeedLimited(start, stop time.Duration) {
	e.speedLimit = rate2.NewRandSpeedLimiter(start, stop)
}

func (e *BaseEngine[KEY, T, W]) Cancel() {
	log.Println("任务取消")
	e.cancel()
	synci.WaitGroupStopWait(&e.wg)

}

func (e *BaseEngine[KEY, T, W]) Run(tasks ...*BaseTask[KEY, T]) {
	e.addWorker()
	e.isRunning = true
	go func() {
		timer := time.NewTimer(time.Second * 5)
		defer timer.Stop()
		var emptyTimes uint
		var readyWorkerCh chan *BaseTask[KEY, T]
		var readyTask *BaseTask[KEY, T]
	loop:
		for {
			if e.workerList.Size > 0 && len(e.taskList) > 0 {
				if readyWorkerCh == nil {
					readyWorkerCh = e.workerList.Pop().taskCh
				}
				if readyTask == nil {
					readyTask = e.taskList.Pop()
				}
			}

			if len(e.taskList) >= int(e.limitWaitTaskCount) {
				select {
				case readyWorker := <-e.workerChan:
					e.workerList.Push(readyWorker)
				case readyWorkerCh <- readyTask:
					readyWorkerCh = nil
					readyTask = nil
				case <-e.ctx.Done():
					break loop
				}
			} else {
				select {
				case readyTaskTmp := <-e.taskChan:
					e.taskList.Push(readyTaskTmp)
				case readyWorker := <-e.workerChan:
					e.workerList.Push(readyWorker)
				case readyWorkerCh <- readyTask:
					readyWorkerCh = nil
					readyTask = nil
				case <-timer.C:
					//检测任务是否已空
					if e.workerList.Size == uint(e.currentWorkerCount) && len(e.taskList) == 0 {
						emptyTimes++
						if emptyTimes > 2 {
							log.Println("task is empty,任务即将结束")
							e.wg.Done()
							break loop
						}
					}
					timer.Reset(time.Second)
				case <-e.ctx.Done():
					break loop
				}
			}
		}
	}()

	e.taskTotalCount += uint64(len(tasks))
	e.wg.Add(len(tasks) + 1)
	for _, task := range tasks {
		if task != nil {
			task.id = gen.GenOrderID()
			e.taskChan <- task
		}
	}

	e.wg.Wait()
	e.isRunning = false
	e.isFinished = true
	log.Println("任务结束")
}

func (e *BaseEngine[KEY, T, W]) newWorker(readyTask *BaseTask[KEY, T]) {
	atomic.AddUint64(&e.currentWorkerCount, 1)
	//id := c.currentWorkerCount
	taskChan := make(chan *BaseTask[KEY, T])
	worker := &Worker[KEY, T, W]{Id: uint(e.currentWorkerCount), taskCh: taskChan}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				log.Println(string(debug.Stack()))
				spew.Dump(readyTask)
				e.wg.Done()
				// 创建一个新的
				e.newWorker(nil)
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
					if e.speedLimit != nil {
						<-e.speedLimit.C
						e.speedLimit.Reset()
					}
					readyTask.BaseTaskFunc(e.ctx)
					atomic.AddUint64(&e.taskDoneCount, 1)
				}
				e.wg.Done()
			case <-e.ctx.Done():
				return
			}
		}
	}()
	e.workers = append(e.workers, worker)
}

func (e *BaseEngine[KEY, T, W]) addWorker() {
	if e.currentWorkerCount != 0 {
		return
	}
	e.newWorker(nil)
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
	task.id = gen.GenOrderID()
	e.taskChan <- task
}

func (e *BaseEngine[KEY, T, W]) AddTasks(tasks ...*BaseTask[KEY, T]) {
	e.wg.Add(len(tasks))
	for _, task := range tasks {
		task.id = gen.GenOrderID()
		e.taskChan <- task
	}
}

func (e *BaseEngine[KEY, T, W]) AddWorker(num int) {
	atomic.AddUint64(&e.limitWorkerCount, uint64(num))
}

func (e *BaseEngine[KEY, T, W]) NewFixedWorker(interval time.Duration) int {
	ch := make(chan *BaseTask[KEY, T])
	e.fixedWorker = append(e.fixedWorker, ch)
	go func() {
		var timer *time.Ticker
		if interval > 0 {
			timer = time.NewTicker(interval)
		}
		for task := range ch {
			if interval > 0 {
				<-timer.C
			}
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
	task.id = gen.GenOrderID()
	go func() {
		ch <- task
	}()
}

func (e *BaseEngine[KEY, T, W]) SyncRun(tasks ...*BaseTask[KEY, T]) {
	panic("TODO")
}

func (e *BaseEngine[KEY, T, W]) RunSingleWorker(tasks ...*BaseTask[KEY, T]) {
	e.NewFixedWorker(0)
	for _, task := range tasks {
		e.AddFixedTask(0, task)
	}
}

func (e *BaseEngine[KEY, T, W]) Release() {
	e.cancel()
	close(e.workerChan)
	close(e.taskChan)
	for _, ch := range e.fixedWorker {
		close(ch)
	}
	if e.speedLimit != nil {
		e.speedLimit.Stop()
	}
	if e.monitor != nil {
		e.monitor.Stop()
	}
}
