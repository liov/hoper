package main

import (
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"os"
)

func main() {
	libDir, _ := osi.CMD("go list -m -f {{.Dir}}  github.com/actliboy/hoper/server/go/lib")
	os.Chdir(libDir)
	osi.CMD("go run " + libDir + "/tools/install.go")
}
