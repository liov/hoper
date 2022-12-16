package main

import (
	"context"
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/initialize/elastic/v8"
	initpostgres "github.com/liov/hoper/server/go/lib/initialize/gormdb/postgres"
	"tools/clawer/timepill"
	"tools/clawer/timepill/es8"
)

type dao1 struct {
	Hoper initpostgres.DB
	Es8   v8.Es `init:"config:elasticsearch8"`
}

func (dao *dao1) Init() {
}

func (dao *dao1) Close() {
}

type dao2 struct {
	Hoper initpostgres.DB
	Es8   v8.Esv2 `init:"config:elasticsearch8"`
}

func (dao *dao2) Init() {
}

func (dao *dao2) Close() {
}

func main() {
	var vdao dao2
	defer initialize.Start(&timepill.Conf, &vdao)()
	timepill.Dao.Hoper = vdao.Hoper
	es8.NewEsDao(context.Background(), vdao.Es8.Client).LoadEs8()
}
