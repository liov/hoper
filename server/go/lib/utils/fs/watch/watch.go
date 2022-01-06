package watch

import (
	"path/filepath"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/fsnotify/fsnotify"
)

type Watch struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}
	handler  Handler
}

type Handler map[string][]func()

func New(interval time.Duration) (*Watch, error) {
	watcher, err := fsnotify.NewWatcher()
	w := &Watch{
		Watcher:  watcher,
		interval: interval,
		done:     make(chan struct{}, 1),
		//1.map和数组做取舍
		handler: make(map[string][]func()),
		//handler:  make(map[string]map[fsnotify.Op]func()),
		//2.提高时间复杂度，用event做key，然后每次事件循环取值
		//handler:  make(map[fsnotify.Event]func()),
	}

	if err == nil {
		go w.run()
	}

	return w, err
}

func (w *Watch) Add(name string, op fsnotify.Op, callback func()) error {
	handle, ok := w.handler[filepath.Clean(name)]
	if !ok {
		err := w.Watcher.Add(name)
		if err != nil {
			return err
		}
		w.handler[filepath.Clean(name)] = make([]func(), 5, 5)
		handle = w.handler[filepath.Clean(name)]
	}
	handle[op-1] = callback
	return nil
}

func (w *Watch) run() {
	var last time.Time
	ev := &fsnotify.Event{}
OuterLoop:
	for {
		select {
		case event, ok := <-w.Watcher.Events:
			if !ok {
				return
			}
			log.Info("event:", event)
			now := time.Now()
			if now.Sub(last) < w.interval && event == *ev {
				continue
			}
			last = now
			ev = &event
			if handle, ok := w.handler[event.Name]; ok {
				for i := range handle {
					if event.Op&fsnotify.Op(i+1) == fsnotify.Op(i+1) && handle[i] != nil {
						handle[i]()
					}
				}
			}
		case err, ok := <-w.Watcher.Errors:
			if !ok {
				return
			}
			log.Error("error:", err)
		case <-w.done:
			break OuterLoop
		}
	}
	close(w.done)
}

func (w *Watch) Close() {
	w.done <- struct{}{}
	w.Watcher.Close()
}
