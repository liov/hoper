package conctrl

import "github.com/actliboy/hoper/server/go/lib/utils/gen"

// todo: implement 协程控制器
type Controller struct {
	id uint64
	ch chan struct{}
}

func NewController(cap int) *Controller {
	return &Controller{
		id: gen.GenOrderID(),
		ch: make(chan struct{}, cap),
	}
}
