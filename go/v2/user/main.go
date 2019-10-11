package main

import (
	"flag"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/initialize/dao"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/utils/log"
)

func main() {
	flag.Parse()
	defer log.Sync()
	defer dao.Dao.Close()
	initialize.Start(config.Conf,config.CustomInit)
	log.Info(*dao.Dao)
}
