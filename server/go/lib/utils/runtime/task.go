package runtimei

import "sync"

type TaskPool struct {
	wg    *sync.WaitGroup
	queue chan struct{}
}

func NewTaskPool(size int) *TaskPool {
	if size <= 0 {
		size = 1
	}
	return &TaskPool{
		queue: make(chan struct{}, size),
		wg:    &sync.WaitGroup{},
	}
}

func (t *TaskPool) Exec(f func()) {
	t.queue <- struct{}{}
	t.wg.Add(1)
	go func() {
		f()
		<-t.queue
		t.wg.Done()
	}()
}

func (t *TaskPool) Wait() {
	t.wg.Wait()
}

// Add 新增一个执行
func (p *TaskPool) Add(delta int) {
	// delta为正数就添加
	for i := 0; i < delta; i++ {
		p.queue <- struct{}{}
	}
	// delta为负数就减少
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

// Done 执行完成减一
func (p *TaskPool) Done() {
	<-p.queue
	p.wg.Done()
}
