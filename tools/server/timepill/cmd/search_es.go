package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/elastic/v8"
	"tools/timepill"
	"tools/timepill/es8"
)

type dao1 struct {
	Hoper db.DB
	Es8   v8.Es `init:"config:elasticsearch8"`
}

func (dao *dao1) Init() {
}

func (dao *dao1) Close() {
}

type dao2 struct {
	Hoper db.DB
	Es8   v8.Esv2 `init:"config:elasticsearch8"`
}

func (dao *dao2) Init() {
}

func (dao *dao2) Close() {
}

func main() {
	var dao dao2
	defer initialize.Start(&timepill.Conf, &dao)()
	es8.NewEsDao(context.Background(), dao.Es8.Client).LoadEs8()
}
