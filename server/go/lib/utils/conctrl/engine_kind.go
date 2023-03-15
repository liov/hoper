package conctrl

import (
	"context"
	"sync"
	"time"
)

// TODO
type KindEngine[KEY comparable, T, W any] struct {
	ctx                           context.Context
	cancel                        context.CancelFunc
	wg                            sync.WaitGroup
	taskDoneCount, taskTotalCount uint64
	kindConfigs                   []*KindConfig[KEY, T, W]
}

type KindConfig[KEY comparable, T, W any] struct {
	limitWorkerCount, currentWorkerCount uint64
	workerChan                           chan *Worker[KEY, T, W]
	taskChan                             chan *Task[KEY, T]
	taskDoneCount, taskTotalCount        uint64
	averageTimeCost                      time.Duration
}

func (e *KindEngine[KEY, T, W]) Run() {

}
