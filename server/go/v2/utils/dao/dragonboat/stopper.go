package tdragonboat

import "sync"

// Stopper is a manager struct for managing worker goroutines.
type Stopper struct {
	shouldStopC chan struct{}
	wg          sync.WaitGroup
}

// NewStopper return a new Stopper instance.
func NewStopper() *Stopper {
	s := &Stopper{
		shouldStopC: make(chan struct{}),
	}

	return s
}

// RunWorker creates a new goroutine and invoke the f func in that new
// worker goroutine.
func (s *Stopper) RunWorker(f func()) {
	s.wg.Add(1)

	go func() {
		f()
		s.wg.Done()
	}()
}

// ShouldStop returns a chan struct{} used for indicating whether the
// Stop() function has been called on Stopper.
func (s *Stopper) ShouldStop() chan struct{} {
	return s.shouldStopC
}

// Wait waits on the internal sync.WaitGroup. It only return when all
// managed worker goroutines are ready to return and called
// sync.WaitGroup.Done() on the internal sync.WaitGroup.
func (s *Stopper) Wait() {
	s.wg.Wait()
}

// Stop signals all managed worker goroutines to stop and wait for them
// to actually stop.
func (s *Stopper) Stop() {
	close(s.shouldStopC)
	s.wg.Wait()
}

// CloseDao closes the internal shouldStopc chan struct{} to signal all
// worker goroutines that they should stop.
func (s *Stopper) Close() {
	close(s.shouldStopC)
}
