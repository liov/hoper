package main

import (
	osi "github.com/hopeio/pandora/utils/os"
	"os"
)

func main() {
	libDir, _ := osi.CMD("go list -m -f {{.Dir}}  github.com/hopeio/pandora")
	os.Chdir(libDir)
	osi.CMD("go run " + libDir + "/tools/install.go")
}
