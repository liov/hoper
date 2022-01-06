package execi

import (
	"log"
	"os"
	"os/exec"

	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
)

func Run(arg string) {
	words := osi.Split(arg)
	cmd := exec.Command(words[0], words[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(cmd.String())
	cmd.Run()
}
