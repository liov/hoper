package conctrl

import (
	"time"
)

type Worker[KEY comparable, T, W any] struct {
	Id     uint
	Kind   Kind
	taskCh chan *BaseTask[KEY, T]
	Props  W
}

type WorkStatistics struct {
	averageTimeCost               time.Duration
	taskDoneCount, taskTotalCount uint64
}
