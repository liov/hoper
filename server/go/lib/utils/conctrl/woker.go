package conctrl

import "time"

type Worker struct {
	Id     uint
	Kind   Kind
	taskCh chan *BaseTask
	WorkStatistics
}

type WorkStatistics struct {
	averageTimeCost               time.Duration
	taskDoneCount, taskTotalCount uint64
}
