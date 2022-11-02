package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/slices"
	"log"
	"sync"
	"time"
)

type TaskInfo struct {
	TaskMeta
	Key      string
	errTimes int
}

type ErrHandleTask interface {
	HasRequestInfo() *TaskInfo
	HasTaskFunc(ctx context.Context) ([]ErrHandleTask, error)
}

type ErrHandleTasks struct {
	tasks      []ErrHandleTask
	generation int
}

type HandleFunc func(ctx context.Context) ([]ErrHandleTask, error)

type Engine struct {
	*BaseEngine
	done        sync.Map
	ReqsChan    chan ErrHandleTasks
	kindHandler []*KindHandler
}

type KindHandler struct {
	Skip bool
	*time.Ticker
	// TODO 指定Kind的Handler
	HandleFun HandleFunc
}

func New(workerCount uint) *Engine {
	return &Engine{
		ReqsChan:   make(chan ErrHandleTasks),
		BaseEngine: NewEngine(workerCount),
	}
}

func (e *Engine) SkipKind(kinds ...Kind) *Engine {
	length := slices.Max(kinds) + 1
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

func (e *Engine) Run(reqs ...ErrHandleTask) {

	go func() {
		for reqs := range e.ReqsChan {
			for _, req := range reqs.tasks {
				if req != nil {
					req.HasRequestInfo().Priority = reqs.generation
					e.BaseEngine.AddTask(e.NewTask(req))
				}
			}
		}
	}()
	tasks := make([]*Task, 0, len(reqs))
	for _, req := range reqs {
		tasks = append(tasks, e.NewTask(req))
	}
	e.BaseEngine.Run(tasks...)
}

func (e *Engine) NewTask(req ErrHandleTask) *Task {

	if req == nil {
		return nil
	}
	reqInfo := req.HasRequestInfo()
	if reqInfo == nil {
		return nil
	}

	var kindHandler *KindHandler
	if e.kindHandler != nil && int(reqInfo.Kind) < len(e.kindHandler) {
		kindHandler = e.kindHandler[reqInfo.Kind]
	}

	if kindHandler != nil && kindHandler.Skip {
		return nil
	}
	if reqInfo.Key != "" {
		if _, ok := e.done.Load(reqInfo.Key); ok {
			return nil
		}
	}
	return &Task{
		TaskMeta: TaskMeta{Kind: reqInfo.Kind},
		Do: func(ctx context.Context) {
			if kindHandler != nil && kindHandler.Ticker != nil {
				<-kindHandler.Ticker.C
			}
			reqs, err := req.HasTaskFunc(ctx)
			if err != nil {
				reqInfo.errTimes++
				log.Println("执行失败", err)
				log.Println("重新执行,key :", reqInfo.Key)
				if reqInfo.errTimes < 5 {
					e.ReqsChan <- ErrHandleTasks{[]ErrHandleTask{req}, req.HasRequestInfo().Priority + 1}
				}
				return
			}
			if reqInfo.Key != "" {
				e.done.Store(reqInfo.Key, struct{}{})
			}
			if len(reqs) > 0 {
				e.ReqsChan <- ErrHandleTasks{reqs, req.HasRequestInfo().Priority + 1}
			}
			return
		},
	}
}
