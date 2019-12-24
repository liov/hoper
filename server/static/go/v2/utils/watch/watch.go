package watch

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/liov/hoper/go/v2/utils/log"
)

type Watch struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}
	callback Callback
}

type Callback map[fsnotify.Event]func()

func New(interval time.Duration) (*Watch, error) {
	watcher, err := fsnotify.NewWatcher()
	w := &Watch{
		Watcher:  watcher,
		interval: interval,
		done:     make(chan struct{}, 1),
		callback: make(map[fsnotify.Event]func()),
	}

	if err == nil {
		go w.run()
	}

	return w, err
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
			now := time.Now()
			if now.Sub(last) < w.interval && event == *ev {
				continue
			}
			last = now
			ev = &event
			if callback, ok := w.callback[event]; ok {
				callback()
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
