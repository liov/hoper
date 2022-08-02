package crawler

import (
	"github.com/actliboy/hoper/server/go/lib/utils/generics/slices"
	"sync"
	"time"
)

type Fail[T any, I comparable] struct {
	Id       I
	FailChan chan T
	CallBack Callback[T]
}

func NewFail[T any, I comparable](id I, callback Callback[T]) *Fail[T, I] {
	return &Fail[T, I]{
		Id:       id,
		FailChan: make(chan T),
		CallBack: callback,
	}
}

func (f *Fail[T, I]) Do(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for txt := range f.FailChan {
			f.CallBack(txt)
		}
		wg.Done()
	}()
}

const defaultLimit = 50

/*
这么写报错

	type TaskCtrl[T any, F generics.Callback[T]] struct {
		wg *sync.WaitGroup
		slices.Index[chan struct{}, string]
		fails []*Fail[T, F]
	}

	func (s *TaskCtrl[T, F]) Add(topic string, i int) {
		s.wg.Add(i)
		// Cannot use 'topic' (type string) as the type O (F)
		c := s.Get(topic)

		c <- struct{}{}
	}
*/
type TaskCtrl[T any, I comparable] struct {
	wg    *sync.WaitGroup
	speed *slices.Index[chan struct{}, string]
	fails []*Fail[T, I]
	timer *time.Timer
}

func (s *TaskCtrl[T, I]) AddTask(topic string, limit int) chan struct{} {
	s.wg.Add(limit)
	ch := make(chan struct{}, limit)
	s.speed.Add(topic, ch)
	return ch
}

func (s *TaskCtrl[T, I]) TaskDo(topic string, f func()) {
	ch := s.speed.Get(topic)
	if ch == nil {
		ch = s.AddTask(topic, defaultLimit)
	}
	s.Add(topic, 1)
	go func() {
		defer s.Done(topic)
		f()
	}()
}

func (s *TaskCtrl[T, I]) Add(topic string, i int) {
	ch := s.speed.Get(topic)
	if ch == nil {
		ch = s.AddTask(topic, i)
	}
	s.wg.Add(i)
	for j := 0; j < i; j++ {
		ch <- struct{}{}
	}
}

func (s *TaskCtrl[T, I]) Done(topic string) {
	ch := s.speed.Get(topic)
	if ch == nil {
		return
	}
	s.wg.Done()
	<-ch
}

func (s *TaskCtrl[T, I]) Wait() {
	s.wg.Wait()
}

func (s *TaskCtrl[T, I]) AddFailHandler(id I, callback Callback[T]) {
	fail := NewFail[T, I](id, callback)
	fail.Do(s.wg)
	s.fails = append(s.fails, fail)
}

func (s *TaskCtrl[T, I]) RemoveFailHandler(id I) {
	for i, f := range s.fails {
		if f.Id == id {
			s.fails = append(s.fails[:i], s.fails[i+1:]...)
			close(f.FailChan)
			break
		}
	}
}

func NewSpeed[T any, I comparable](speed time.Duration) *TaskCtrl[T, I] {
	return &TaskCtrl[T, I]{
		wg:    new(sync.WaitGroup),
		speed: slices.NewIndex[chan struct{}, string](),
		fails: make([]*Fail[T, I], 0),
		timer: time.NewTimer(speed),
	}
}
