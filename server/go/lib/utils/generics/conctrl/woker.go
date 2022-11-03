package conctrl

import (
	"time"
)

type Worker[T, W any] struct {
	Id     uint
	Kind   Kind
	taskCh chan *BaseTask[T]
	Props  W
}

type WorkStatistics struct {
	averageTimeCost               time.Duration
	taskDoneCount, taskTotalCount uint64
}
