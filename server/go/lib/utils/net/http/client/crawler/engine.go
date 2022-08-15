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
	ctrlEngine   *conctrl.Engine
	visited      sync.Map
	reqsChan     chan []*Request
	excludeKinds []bool
	timer        []*time.Ticker
}

func New(workerCount uint) *Engine {
	return &Engine{
		reqsChan:   make(chan []*Request),
		ctrlEngine: conctrl.NewEngine(workerCount),
	}
}

func (e *Engine) SkipKind(kinds ...conctrl.Kind) *Engine {
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

func (e *Engine) Timer(kind conctrl.Kind, interval time.Duration) *Engine {
	if e.timer == nil {
		e.timer = make([]*time.Ticker, int(kind)+1)
	}
	if int(kind)+1 > len(e.timer) {
		e.timer = append(e.timer, make([]*time.Ticker, int(kind)+1-len(e.timer))...)
	}
	e.timer[kind] = time.NewTicker(interval)
	return e
}

func (e *Engine) Run(reqs ...*Request) {

	go func() {
		for reqs := range e.reqsChan {
			for _, req := range reqs {
				if e.excludeKinds != nil && int(req.Kind) < len(e.excludeKinds) && e.excludeKinds[req.Kind] {
					continue
				}
				e.ctrlEngine.AddTask(e.NewTask(req))
			}
		}
	}()
	tasks := make([]*conctrl.Task, 0, len(reqs))
	for _, req := range reqs {
		tasks = append(tasks, e.NewTask(req))
	}
	e.ctrlEngine.Run(tasks...)
}

func (e *Engine) NewTask(req *Request) *conctrl.Task {
	if _, ok := e.visited.Load(req.Url); ok {
		return nil
	}
	return &conctrl.Task{
		Kind: req.Kind,
		Do: func(ctx context.Context) {
			if e.timer != nil && int(req.Kind) < len(e.timer) && e.timer[req.Kind] != nil {
				<-e.timer[req.Kind].C
			}
			if _, ok := e.visited.Load(req.Url); ok {
				return
			}
			reqs, err := req.HandleFun(ctx, req.Url)
			if err != nil {
				log.Println("爬取失败", err)
				log.Println("重新爬取,url :", req.Url)
				e.reqsChan <- []*Request{req}
				return
			}
			e.visited.Store(req.Url, struct{}{})
			if len(reqs) > 0 {
				e.reqsChan <- reqs
			}
			return
		},
	}
}
