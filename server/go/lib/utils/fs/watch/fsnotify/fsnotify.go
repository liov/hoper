package fsnotify

import (
	"path/filepath"
	"time"

	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/fsnotify/fsnotify"
)

type Watch struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}
	handler  map[string]Callback
}

type Callback struct {
	lastModTime time.Time
	callbacks   []func()
}

type Handler map[string]Callback

func New(interval time.Duration) (*Watch, error) {
	watcher, err := fsnotify.NewWatcher()
	w := &Watch{
		Watcher:  watcher,
		interval: interval,
		done:     make(chan struct{}, 1),
		//1.map和数组做取舍
		handler: make(map[string]Callback),
		//Handler:  make(map[string]map[fsnotify.Op]func()),
		//2.提高时间复杂度，用event做key，然后每次事件循环取值
		//Handler:  make(map[fsnotify.Event]func()),
	}

	if err == nil {
		go w.run()
	}

	return w, err
}

func (w *Watch) Add(name string, op fsnotify.Op, callback func()) error {
	name = filepath.Clean(name)
	handle, ok := w.handler[name]
	if !ok {
		err := w.Watcher.Add(name)
		if err != nil {
			return err
		}
		w.handler[name] = Callback{
			callbacks: make([]func(), 5, 5),
		}
		handle = w.handler[name]
	}
	handle.callbacks[op-1] = callback
	return nil
}

func (w *Watch) run() {
	ev := &fsnotify.Event{}
OuterLoop:
	for {
		select {
		case event, ok := <-w.Watcher.Events:
			if !ok {
				return
			}
			log.Info("event:", event)
			ev = &event
			if handle, ok := w.handler[event.Name]; ok {
				now := time.Now()
				if now.Sub(handle.lastModTime) < w.interval && event == *ev {
					continue
				}
				handle.lastModTime = now
				for i := range handle.callbacks {
					if event.Op&fsnotify.Op(i+1) == fsnotify.Op(i+1) && handle.callbacks[i] != nil {
						handle.callbacks[i]()
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
}

func (w *Watch) Close() {
	w.done <- struct{}{}
	close(w.done)
	w.Watcher.Close()
}
