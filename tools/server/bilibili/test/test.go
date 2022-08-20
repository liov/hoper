package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	initpostgres "github.com/actliboy/hoper/server/go/lib/initialize/db/postgres"
	"tools/bilibili/config"
)

type dao3 struct {
	Hoper initpostgres.DB
}

func (d dao3) Init() {
}

func main() {
	dao := &dao3{}
	defer initialize.Start(config.Conf, dao)()
}
