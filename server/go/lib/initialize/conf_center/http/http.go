package http

import (
	http_fs "github.com/actliboy/hoper/server/go/lib/utils/net/http/fs"
	http_fs_watch "github.com/actliboy/hoper/server/go/lib/utils/net/http/fs/watch"
	"time"
)

type Config struct {
	Interval int64
	Url      string
}

// 本地配置
func (cc *Config) HandleConfig(handle func([]byte)) error {

	if cc.Interval == 0 {
		cc.Interval = 5
	}

	watch := http_fs_watch.New(time.Second * time.Duration(cc.Interval))

	callback := func(hfile *http_fs.FileInfo) {
		handle(hfile.Binary)
	}

	return watch.Add(cc.Url, callback)
}