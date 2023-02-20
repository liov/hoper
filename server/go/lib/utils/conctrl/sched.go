package conctrl

import (
	"context"
)

type SchedulingContext struct {
	context.Context
	run  chan struct{}
	Done bool
}

func (c *SchedulingContext) Monitor() func() {
	c.run <- struct{}{}
	return func() {
		<-c.run
		if len(c.run) == 0 {
			c.Done = true
		}
	}
}
