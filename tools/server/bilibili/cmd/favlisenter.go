package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"tools/bilibili/config"
)

func main() {

	defer initialize.Start(config.Conf, nil)()

}
