package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/slices"
	"log"
	"sync"
	"time"
)

type Engine struct {
	*conctrl.Engine
	visited     sync.Map
	ReqsChan    chan []*Request
	kindHandler []*KindHandler
}

type KindHandler struct {
	Skip bool
	*time.Ticker
	// TODO 指定Kind的Handler
	HandleFun HandleFun
}

func New(workerCount uint) *Engine {
	return &Engine{
		ReqsChan: make(chan []*Request),
		Engine:   conctrl.NewEngine(workerCount),
	}
}

func (e *Engine) SkipKind(kinds ...conctrl.Kind) *Engine {
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

func (e *Engine) Timer(kind conctrl.Kind, interval time.Duration) *Engine {
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

func (e *Engine) Run(reqs ...*Request) {

	go func() {
		for reqs := range e.ReqsChan {
			for _, req := range reqs {
				e.Engine.AddTask(e.NewTask(req))
			}
		}
	}()
	tasks := make([]*conctrl.Task, 0, len(reqs))
	for _, req := range reqs {
		tasks = append(tasks, e.NewTask(req))
	}
	e.Engine.Run(tasks...)
}

func (e *Engine) NewTask(req *Request) *conctrl.Task {
	if req == nil || req.TaskFun == nil {
		return nil
	}
	var kindHandler *KindHandler
	if e.kindHandler != nil && int(req.Kind) < len(e.kindHandler) {
		kindHandler = e.kindHandler[req.Kind]
	}

	if kindHandler != nil && kindHandler.Skip {
		return nil
	}
	if req.Key != "" {
		if _, ok := e.visited.Load(req.Key); ok {
			return nil
		}
	}
	return &conctrl.Task{
		TaskMeta: conctrl.TaskMeta{Kind: req.Kind},
		Do: func(ctx context.Context) {
			if kindHandler != nil && kindHandler.Ticker != nil {
				<-kindHandler.Ticker.C
			}
			reqs, err := req.TaskFun(ctx)
			if err != nil {
				req.errTimes++
				log.Println("爬取失败", err)
				log.Println("重新爬取,url :", req.Key)
				if req.errTimes < 5 {
					e.ReqsChan <- []*Request{req}
				}
				return
			}
			if req.Key != "" {
				e.visited.Store(req.Key, struct{}{})
			}
			if len(reqs) > 0 {
				e.ReqsChan <- reqs
			}
			return
		},
	}
}
