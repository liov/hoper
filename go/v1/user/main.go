package main

import (
	"github.com/liov/hoper/go/v1/initialize"
	"github.com/liov/hoper/go/v1/initialize/dao"
	"github.com/liov/hoper/go/v1/user/internal/config"

	"github.com/liov/hoper/go/v1/utils/log"
)

func main() {
	defer log.Sync()
	defer dao.Dao.Close()
	initialize.Start(config.Conf)
	log.Info(*dao.Dao)
}
