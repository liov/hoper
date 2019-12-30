package main

import (
	"os/signal"
	"syscall"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/user/internal/server"
)

func main() {
	defer initialize.Start(config.Conf, dao.Dao)()
Loop:
	for {
		signal.Notify(server.SignalChan(),
			// kill -SIGINT XXXX 或 Ctrl+c
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-server.SignalChan():
			break Loop
		default:
			server.Serve()
		}
	}
}
