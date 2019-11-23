package main

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type foo struct {
	time string `json:"time"`
}

func TestSql(t *testing.T) {
	defer sqlDB.Close()
	in := `INSERT INTO test (time2) VALUES (?);`
	f:= foo{time:"2019-10-31 15:06:00.772"}
	_,err:=sqlDB.Exec(in,f.time)
	if err!= nil{
		log.Println(err)
	}
}
