package main

import (
	"github.com/liov/hoper/go/initialize"
	"github.com/liov/hoper/go/user/internal/config"
	"github.com/liov/hoper/go/user/internal/dao"
	"github.com/liov/hoper/go/utls/log"
)

func main() {
	defer log.Sync()
	initialize.Start(config.Conf)
	log.Info(dao.Dao)
}
