package conctrl

import (
	"context"
)

type RequestInfo struct {
	TaskMeta
	Key      string
	errTimes int
}

type Request interface {
	RequestInfo() *RequestInfo
	TaskFunc(ctx context.Context) ([]Request, error)
}

type Requests struct {
	reqs       []Request
	generation int
}
