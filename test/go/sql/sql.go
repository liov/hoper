package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init(){
	user:=flag.String( "u", "", "user")
	ip:=flag.String( "ip", "", "ip")
	flag.Parse()
	var err error
	db, err = sql.Open("mysql", *user+":123456@tcp("+*ip+":3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
}

type foo struct {
	time string `json:"time"`
}

func main() {
	defer db.Close()
	in := `INSERT INTO test (time2) VALUES (?);`
	f:= foo{time:"2019-10-31 15:06:00.772"}
	_,err:=db.Exec(in,f.time)
	if err!= nil{
		log.Println(err)
	}
}
