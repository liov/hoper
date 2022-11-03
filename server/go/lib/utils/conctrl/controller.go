package conctrl

type Controller chan func() error

func (c Controller) AddTask(f func() error) {
	go func() {
		c <- f
	}()
}

func (c Controller) Start() {
	for f := range c {
		err := f()
		if err != nil {
			c.AddTask(f)
		}
	}
}
