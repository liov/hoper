package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"tools/clawer/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	timepill.CreateTable()
}
