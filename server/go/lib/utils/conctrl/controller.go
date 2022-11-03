package conctrl

type Controller struct {
	ch chan func() error
}

func (c *Controller) AddTask(f func() error) {
	go func() {
		c.ch <- f
	}()
}

func (c *Controller) Start() {
	for f := range c.ch {
		err := f()
		if err != nil {
			c.AddTask(f)
		}
	}
}
