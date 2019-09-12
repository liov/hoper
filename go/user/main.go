package main

import (
	"github.com/liov/hoper/go/user/initialize"
	"github.com/liov/hoper/go/user/internal/dao"
	"github.com/liov/hoper/go/utls/log"
)

func main() {
	defer log.Sync()
	initialize.Start()
	log.Info(dao.Dao)
}
