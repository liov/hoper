package poll

import (
	"io/fs"
	"time"
)

type Watch struct {
	interval time.Duration
	done     chan struct{}
	handler  map[string]Callback
}

type Callback struct {
	lastModTime time.Time
	callbacks   []func(file fs.FileInfo)
}

type Handler map[string]Callback
