package main

import (
	"context"
	"fmt"
	"github.com/liov/hoper/server/go/lib/initialize"
	initpostgres "github.com/liov/hoper/server/go/lib/initialize/gormdb/postgres"
	"tools/clawer/bilibili/config"
	"tools/clawer/bilibili/dao"
)

type dao1 struct {
	Hoper initpostgres.DB
}

func (d dao1) Init() {
}

func main() {
	var daod dao1
	defer initialize.Start(config.Conf, &daod)()
	//daod.Hoper.Migrator().CreateTable(&dao.View{}, &dao.Video{}, &dao.ViewBak{})
	fmt.Println(dao.NewDao(context.Background(), daod.Hoper.DB).LastCreated(dao.TableNameView))
	fmt.Println(dao.NewDao(context.Background(), daod.Hoper.DB).LastRecordAid())
}
