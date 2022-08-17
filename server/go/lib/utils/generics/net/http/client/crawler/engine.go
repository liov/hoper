package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"log"
	"sync"
	"time"
)

type Engine struct {
	ctrlEngine   *conctrl.Engine
	visited      sync.Map
	reqsChan     chan []*crawler.Request
	kindHandlers []*KindHandler
}

type KindHandler struct {
	conctrl.KindHandler
	handleFun crawler.HandleFun
	timer     *time.Ticker
}

func New(workerCount uint) *Engine {
	return &Engine{
		reqsChan:   make(chan []*crawler.Request),
		ctrlEngine: conctrl.NewEngine(workerCount),
	}
}

func (e *KindHandler) HandleFun(handleFun crawler.HandleFun) *KindHandler {
	e.handleFun = handleFun
	return e
}

func (e *KindHandler) Timer(interval time.Duration) *KindHandler {
	e.timer = time.NewTicker(interval)
	return e
}

func (e *Engine) Run(reqs ...*crawler.Request) {

	go func() {
		for reqs := range e.reqsChan {
			for _, req := range reqs {
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

func (e *Engine) NewTask(req *crawler.Request) *conctrl.Task {
	if _, ok := e.visited.Load(req.Url); ok {
		return nil
	}
	if e.kindHandlers == nil || int(req.Kind) < len(e.kindHandlers) || e.kindHandlers[req.Kind] == nil {
		return nil
	}
	handler := e.kindHandlers[req.Kind]
	req.HandleFun = handler.handleFun
	return &conctrl.Task{
		Kind: req.Kind,
		Do: func(ctx context.Context) {
			if handler.timer != nil {
				<-handler.timer.C
			}
			if _, ok := e.visited.Load(req.Url); ok {
				return
			}
			reqs, err := req.HandleFun(ctx, req.Url)
			if err != nil {
				log.Println("爬取失败", err)
				log.Println("重新爬取,url :", req.Url)
				e.reqsChan <- reqs
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
