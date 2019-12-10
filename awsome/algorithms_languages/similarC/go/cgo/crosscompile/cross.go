package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/exp/shiny/driver"
	_ "test/cgo/crosscompile/a"
	"test/cgo/crosscompile/excel"
	"test/cgo/crosscompile/logrus"
	"test/cgo/crosscompile/mobile"
)

func main() {
	excel.Excel()
	//unicorn.Unicorn()
	//a.Hello()
	mobile.Mobile()
	logrus.Logrus()
	ctxt := new(build.Context)
	*ctxt = build.Default
	ctxt.CgoEnabled = false

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	aDir := filepath.Join(wd, "a")
	p, err := ctxt.ImportDir(aDir, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Join(p.Imports, "\n"))
}
