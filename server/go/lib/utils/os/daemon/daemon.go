package daemon

import (
	"flag"
	"os"
	"os/exec"

	"github.com/liov/hoper/server/go/lib/utils/log"
)

var d bool

func init() {
	flag.BoolVar(&d, "d", false, "守护进程")
	if !flag.Parsed() {
		flag.Parse()
	}

	if d {
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "-d=true" {
				os.Args[i] = "-d=false"
			}
		}
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		_ = cmd.Start()
		log.Info("[PID]", cmd.Process.Pid)
		os.Exit(0)
	}
}
