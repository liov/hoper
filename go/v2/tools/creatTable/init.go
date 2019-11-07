package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/liov/hoper/go/v2/protobuf/user"
)

var(
	ormDB *gorm.DB
)


func init() {
	var config = struct {
		User string
		PassWord string
		Host string
	}{}
	err := configor.New(&configor.Config{Debug: true}).
		Load(&config, "./add-config.toml")
	if err != nil {
		log.Fatal(err)
	}
	url:= fmt.Sprintf("%s:%s@tcp(%s:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",config.User,config.PassWord,config.Host)

	ormDB, err = gorm.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}
	ormDB.SingularTable(true)
	ormDB.LogMode(true)
}

func main() {
	ormDB.CreateTable(&user.User{})
}
