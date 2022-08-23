package conctrl

import (
	"context"
	"sync"
	"time"
)

// TODO
type KindEngine struct {
	ctx                           context.Context
	cancel                        context.CancelFunc
	wg                            sync.WaitGroup
	taskDoneCount, taskTotalCount uint64
	engines                       []*Engine
}

type kindEngine struct {
	limitWorkerCount, currentWorkerCount uint64
	workerChan                           chan *Worker
	taskChan                             chan *Task
	kindHandler                          KindHandler
	taskDoneCount, taskTotalCount        uint64
	averageTimeCost                      time.Duration
}

func (e *KindEngine) Run() {

}
