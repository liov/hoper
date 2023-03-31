package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hopeio/pandora/initialize"
	"github.com/hopeio/pandora/initialize/gormdb/postgres"
	"log"
	"time"
)

type PConfig struct {
}

func (c *PConfig) Init() {

}

type PDao struct {
	Hoper postgres.DB
}

func (d *PDao) Init() {

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
