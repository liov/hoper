package conctrl

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/slices"
	"log"
	"sync"
	"time"
)

type Engine struct {
	*BaseEngine
	done        sync.Map
	TasksChan   chan TaskInterfaces
	kindHandler []*KindHandler
}

type KindHandler struct {
	Skip bool
	*time.Ticker
	// TODO 指定Kind的Handler
	HandleFun TaskFunc
}

func NewEngine(workerCount uint) *Engine {
	return &Engine{
		TasksChan:  make(chan TaskInterfaces),
		BaseEngine: NewBaseEngine(workerCount),
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

func (e *Engine) Run(reqs ...TaskInterface) {

	go func() {
		for reqs := range e.TasksChan {
			for _, req := range reqs.tasks {
				if req != nil {
					req.HasTask().Priority = reqs.generation
					e.BaseEngine.AddTask(e.NewTask(req))
				}
			}
		}
	}()
	tasks := make([]*BaseTask, 0, len(reqs))
	for _, req := range reqs {
		tasks = append(tasks, e.NewTask(req))
	}
	e.BaseEngine.Run(tasks...)
}

func (e *Engine) NewTask(req TaskInterface) *BaseTask {

	if req == nil {
		return nil
	}
	reqInfo := req.HasTask()
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
	return &BaseTask{
		BaseTaskMeta: BaseTaskMeta{Id: reqInfo.Id, Priority: reqInfo.Priority},
		BaseTaskFunc: func(ctx context.Context) {
			if kindHandler != nil && kindHandler.Ticker != nil {
				<-kindHandler.Ticker.C
			}
			reqs, err := reqInfo.TaskFunc(ctx)
			if err != nil {
				reqInfo.ErrTimes++
				log.Println("执行失败", err)
				log.Println("重新执行,key :", reqInfo.Key)
				if reqInfo.ErrTimes < 5 {
					e.TasksChan <- TaskInterfaces{[]TaskInterface{req}, reqInfo.Priority + 1}
				}
				return
			}
			if reqInfo.Key != "" {
				e.done.Store(reqInfo.Key, struct{}{})
			}
			if len(reqs) > 0 {
				e.TasksChan <- TaskInterfaces{reqs, reqInfo.Priority + 1}
			}
			return
		},
	}
}
