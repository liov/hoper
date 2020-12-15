package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB1        *gorm.DB
	DB2        *gorm.DB
	configPath = "D:\\config/add-config.toml"
)

func init() {
	log.SetFlags(log.Llongfile)
	var config = struct {
		DB1 struct {
			User     string
			PassWord string
			Host     string
			DataBase string
		}
		DB2 struct {
			User     string
			PassWord string
			Host     string
			DataBase string
		}
	}{}
	err := configor.New(&configor.Config{Debug: true}).
		Load(&config, configPath)
	if err != nil {
		log.Fatal(err)
	}
	url1 := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB1.User, config.DB1.PassWord, config.DB1.Host, config.DB1.DataBase)
	url2 := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB2.User, config.DB2.PassWord, config.DB2.Host, config.DB2.DataBase)

	DB1, err = gorm.Open("mysql", url1)
	if err != nil {
		log.Fatalln(err)
	}
	DB2, err = gorm.Open("mysql", url2)
	if err != nil {
		log.Fatalln(err)
	}
	DB1.SingularTable(true)
	DB1.SingularTable(true)
	DB2.LogMode(true)
	DB2.LogMode(true)
}

type RecRule struct {
	Id      uint64
	Filters []Filter `gorm:"foreignkey:RecRuleId"`
}

type Filter struct {
	RecRuleId uint64 `json:"rec_rule_id"`
	Field     string
	Method    string
	Value     string
}

func (*Filter) TableName() string {
	return "res_filter"
}

//牛啊，可以跨数据库preload
func main() {
	var r RecRule
	DB1.CreateTable(&r)
	r.Id = 9
	DB1.Save(&r)
	err := DB1.Preload("Filters", func(db *gorm.DB) *gorm.DB {
		return DB2
	}).Find(&r).Error
	if err != nil {
		log.Println(err)
	}
	log.Println(r)
}
