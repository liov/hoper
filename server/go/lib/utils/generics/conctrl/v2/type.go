package v2

import "github.com/actliboy/hoper/server/go/lib/utils/conctrl"

type Task[T TaskProps] struct {
	Id    uint
	Kind  conctrl.Kind
	Props T
}

type TaskProps interface {
	NewTaskFun(id uint, kind conctrl.Kind) conctrl.BaseTaskFunc
}
