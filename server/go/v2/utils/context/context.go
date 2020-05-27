package context2

import (
	"context"
	"time"
)

var vctx = new(valueContext)

func ValueContext() context.Context {
	return vctx
}

type valueContext struct {
	value interface{}
}

func (*valueContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*valueContext) Done() <-chan struct{} {
	return nil
}

func (*valueContext) Err() error {
	return nil
}

func (ctx *valueContext) Value(key interface{}) interface{} {
	return ctx.value
}

func (ctx *valueContext) SetValue(value interface{}) {
	ctx.value = value
}

type Context struct {
	context.Context
	run  chan struct{}
	Done bool
}

func (c *Context) Monitor() func() {
	c.run <- struct{}{}
	return func() {
		<-c.run
		if len(c.run) == 0 {
			c.Done = true
		}
	}
}
