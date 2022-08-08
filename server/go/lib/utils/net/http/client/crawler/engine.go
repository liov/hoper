package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"log"
	"sync"
)

type Engine struct {
	ctrlEngine *conctrl.Engine
	visited    sync.Map
	reqsChan   chan []*Request
}

func New(workerCount int) *Engine {
	return &Engine{
		reqsChan:   make(chan []*Request),
		ctrlEngine: conctrl.NewEngine(workerCount),
	}
}

func (e *Engine) Run(reqs ...*Request) {
	tasks := make([]conctrl.Task, 0, len(reqs))
	go func() {
		for reqs := range e.reqsChan {
			for _, req := range reqs {
				e.ctrlEngine.AddTask(e.NewTask(req))
			}
		}
	}()
	for _, req := range reqs {
		tasks = append(tasks, e.NewTask(req))
	}
	e.ctrlEngine.Run(tasks...)
}

func (e *Engine) NewTask(req *Request) conctrl.Task {
	if _, ok := e.visited.Load(req.Url); ok {
		return nil
	}
	return func(ctx context.Context) {
		if _, ok := e.visited.Load(req.Url); ok {
			return
		}
		reqs, err := req.HandleFun(req.Url)
		if err != nil {
			//req.FailFun(err)
			log.Println("爬取失败", err)
			log.Println("重新爬取,url :", req.Url)
			e.reqsChan <- []*Request{req}
			return
		}
		e.visited.Store(req.Url, struct{}{})
		if len(reqs) > 0 {
			e.reqsChan <- reqs
		}
	}
}
