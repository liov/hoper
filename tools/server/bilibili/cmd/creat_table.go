package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"tools/bilibili/config"
	"tools/bilibili/model"
)

type dao1 struct {
	Hoper db.DB
}

func (d dao1) Init() {
}

func main() {
	var dao dao1
	defer initialize.Start(config.Conf, &dao)()
	dao.Hoper.Migrator().CreateTable(&model.View{}, &model.Video{})
}
