package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hopeio/cherry/initialize"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/postgres"
	"log"
	"time"
)

type PConfig struct {
	initialize.EmbeddedPresets
}

type PDao struct {
	initialize.EmbeddedPresets
	Hoper postgres.DB
}

var pdao PDao

type Test struct {
	Id        int
	DeletedAt time.Time
}

func main() {
	defer initialize.Start(nil, &pdao)()
	var tests []*Test
	pdao.Hoper.Find(&tests)
	spew.Dump(tests)
	zeroTime := time.Time{}
	for _, test := range tests {
		log.Println(test.DeletedAt == zeroTime)
		log.Println(test.DeletedAt == zeroTime.Local())
		log.Println(test.DeletedAt.UTC() == zeroTime)
	}
	pdao.Hoper.Create(&Test{Id: 7})
}
