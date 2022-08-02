package _go

import (
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"os"
	"strings"
)

const goListDir = `go list -m -f {{.Dir}} `
const GOPATHKey = "GOPATH"

var gopath, modPath string

func init() {
	if gopath == "" {
		gopath = os.Getenv(GOPATHKey)
	}
	if gopath != "" && !strings.HasSuffix(gopath, "/") {
		gopath = gopath + "/"
	}
	modPath = gopath + "pkg/mod/"
}

func GetDepDir(dep string) string {
	if !strings.Contains(dep, "@") {
		return modDepDir(dep)
	}
	depPath := modPath + dep
	_, err := os.Stat(depPath)
	if os.IsNotExist(err) {
		depPath = modDepDir(dep)
	}
	return depPath
}

func modDepDir(dep string) string {
	depPath, _ := osi.CMD(goListDir + dep)
	if depPath == "" {
		osi.CMD("go get " + dep)
		depPath, _ = osi.CMD(goListDir + dep)
	}
	return depPath
}
