package _go

import (
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"os"
	"strings"
)

const goList = `go list -m -f {{.Dir}} `
const GOPATHKey = "GOPATH"

func GetDepDir(gopath, dep string) string {
	//
	if gopath == "" {
		gopath = os.Getenv(GOPATHKey)
	}
	if gopath != "" && !strings.HasSuffix(gopath, "/") {
		gopath = gopath + "/"
	}
	modPath := gopath + "pkg/mod/"
	depPath := modPath + dep
	_, err := os.Stat(depPath)
	if os.IsNotExist(err) {
		osi.CMD("go get " + dep)
		depPath, _ = osi.CMD(
			goList + "github.com/googleapis/googleapis",
		)
	}
	return depPath
}
