package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	sqlDB      *sql.DB
	ormDB      *gorm.DB
	configPath = "D:\\add-config.toml"
)

func init() {
	var config = struct {
		User     string
		PassWord string
		Host     string
	}{}
	err := configor.New(&configor.Config{Debug: true}).
		Load(&config, configPath)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.PassWord, config.Host)
	sqlDB, err = sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	ormDB, err = gorm.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}
	ormDB.SingularTable(true)
	ormDB.LogMode(true)
}
