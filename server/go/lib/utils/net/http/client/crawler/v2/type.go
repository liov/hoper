package crawler

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
)

type RequestInfo struct {
	conctrl.TaskMeta
	Key      string
	errTimes int
}

type Request interface {
	RequestInfo() *RequestInfo
	TaskFunc(ctx context.Context) ([]Request, error)
}
