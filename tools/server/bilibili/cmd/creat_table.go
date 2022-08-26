package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	initpostgres "github.com/actliboy/hoper/server/go/lib/initialize/db/postgres"
	"tools/bilibili/config"
	"tools/bilibili/dao"
)

type dao1 struct {
	Hoper initpostgres.DB
}

func (d dao1) Init() {
}

func main() {
	var daod dao1
	defer initialize.Start(config.Conf, &daod)()
	daod.Hoper.Migrator().CreateTable(&dao.View{}, &dao.Video{}, &dao.ViewBak{})
}
