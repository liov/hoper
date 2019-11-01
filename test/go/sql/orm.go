package main

import (
	"flag"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	user := flag.String("u", "", "user")
	ip := flag.String("ip", "", "ip")
	flag.Parse()
	var err error
	db, err = gorm.Open("mysql", *user+":123456@tcp("+*ip+":3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
}

type User struct {
	ID          uint64 `gorm:"primary_key" json:"id"`
	ActivatedAt string `json:"-"` //激活时间
	Name        string `gorm:"type:varchar(10);not null" json:"name"`
	Password    string `gorm:"type:varchar(100)" json:"-"`
}

func main() {
	user := User{ActivatedAt: "2019-10-31 15:06:00.772", Name: "test", Password: "123"}
	err := db.Save(&user).Error
	if err != nil {
		log.Fatalln(err)
	}
	user = User{}
	db.First(&user)
	log.Println(user)
}
