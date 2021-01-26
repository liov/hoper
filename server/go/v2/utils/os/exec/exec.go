package execi

import (
	"os"
	"os/exec"

	osi "github.com/liov/hoper/go/v2/utils/os"
)

func Run(arg string) {
	words := osi.Split(arg)
	cmd := exec.Command(words[0], words[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
