package conctrl

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/gen"
	"github.com/liov/hoper/server/go/lib/utils/slices"
	"log"
	"strconv"
	"sync"
	"time"
)

type Engine struct {
	*BaseEngine
	done        sync.Map
	kindHandler []*KindHandler
	errHandler  func(task *Task)
	errChan     chan *Task
}

type KindHandler struct {
	Skip bool
	*time.Ticker
	// TODO 指定Kind的Handler
	HandleFun TaskFunc
}

func NewEngine(workerCount uint) *Engine {
	return &Engine{
		BaseEngine: NewBaseEngine(workerCount),
		errHandler: func(task *Task) {
			log.Println("错误次数达到"+strconv.Itoa(task.ErrTimes)+"次,放弃执行:", task.Errs)
		},
		errChan: make(chan *Task),
	}
}

func (e *Engine) SkipKind(kinds ...Kind) *Engine {
	length := slices.Uint8Max(kinds) + 1
	if e.kindHandler == nil {
		e.kindHandler = make([]*KindHandler, length)
	}
	if int(length) > len(e.kindHandler) {
		e.kindHandler = append(e.kindHandler, make([]*KindHandler, int(length)-len(e.kindHandler))...)
	}
	for _, kind := range kinds {
		if e.kindHandler[kind] == nil {
			e.kindHandler[kind] = &KindHandler{Skip: true}
		} else {
			e.kindHandler[kind].Skip = true
		}

	}
	return e
}
func (e *Engine) StopAfter(interval time.Duration) *Engine {
	time.AfterFunc(interval, e.Cancel)
	return e
}

func (e *Engine) ErrHandler(errHandler func(group *Task)) *Engine {
	e.errHandler = errHandler
	return e
}

func (e *Engine) Timer(kind Kind, interval time.Duration) *Engine {
	if e.kindHandler == nil {
		e.kindHandler = make([]*KindHandler, int(kind)+1)
	}
	if int(kind)+1 > len(e.kindHandler) {
		e.kindHandler = append(e.kindHandler, make([]*KindHandler, int(kind)+1-len(e.kindHandler))...)
	}
	if e.kindHandler[kind] == nil {
		e.kindHandler[kind] = &KindHandler{Ticker: time.NewTicker(interval)}
	} else {
		e.kindHandler[kind].Ticker = time.NewTicker(interval)
	}
	return e
}

func (e *Engine) Run(tasks ...TaskInterface) {

	baseTasks := make([]*BaseTask, 0, len(tasks))
	for _, req := range tasks {
		baseTasks = append(baseTasks, e.NewTask(req))
	}
	go func() {
		for group := range e.errChan {
			e.errHandler(group)
		}
	}()
	e.BaseEngine.Run(baseTasks...)
}

func (e *Engine) NewTask(task TaskInterface) *BaseTask {

	if task == nil {
		return nil
	}
	taskInfo := task.HasTask()
	if taskInfo == nil {
		return nil
	}
	taskInfo.Id = gen.GenOrderID()

	var kindHandler *KindHandler
	if e.kindHandler != nil && int(taskInfo.Kind) < len(e.kindHandler) {
		kindHandler = e.kindHandler[taskInfo.Kind]
	}

	if kindHandler != nil && kindHandler.Skip {
		return nil
	}
	if taskInfo.Key != "" {
		if _, ok := e.done.Load(taskInfo.Key); ok {
			return nil
		}
	}
	return &BaseTask{
		BaseTaskMeta: BaseTaskMeta{Id: taskInfo.Id, Priority: taskInfo.Priority},
		BaseTaskFunc: func(ctx context.Context) {
			if kindHandler != nil && kindHandler.Ticker != nil {
				<-kindHandler.Ticker.C
			}
			tasks, err := taskInfo.TaskFunc(ctx)
			if err != nil {
				taskInfo.ErrTimes++
				taskInfo.Errs = append(taskInfo.Errs, err)
				log.Println(taskInfo.Key, "执行失败", err, "重新执行")
				if taskInfo.ErrTimes < 5 {
					e.AsyncAddTask(taskInfo.Priority+1, task)
				}
				if taskInfo.ErrTimes == 5 {
					e.errChan <- taskInfo
				}
				return
			}
			if taskInfo.Key != "" {
				e.done.Store(taskInfo.Key, struct{}{})
			}
			if len(tasks) > 0 {
				e.AsyncAddTask(taskInfo.Priority+1, tasks...)
			}
			return
		},
	}
}

func (e *Engine) AddTasks(generation int, tasks ...TaskInterface) {
	for _, task := range tasks {
		if task != nil {
			task.HasTask().Priority = generation
			e.BaseEngine.AddTask(e.NewTask(task))
		}
	}
}

func (e *Engine) AsyncAddTask(generation int, tasks ...TaskInterface) {
	go func() {
		for _, task := range tasks {
			if task != nil {
				task.HasTask().Priority = generation
				e.BaseEngine.AddTask(e.NewTask(task))
			}
		}
	}()
}
