package main

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)


type User struct {
	ID          uint64 `gorm:"primary_key" json:"id"`
	ActivatedAt string `json:"-"` //激活时间
	Name        string `gorm:"type:varchar(10);not null" json:"name"`
	Password    string `gorm:"type:varchar(100)" json:"-"`
}

func TestORM(t *testing.T) {
	defer ormDB.Close()
	user := User{ActivatedAt: "2019-10-31 15:06:00.772", Name: "test", Password: "123"}
	err := ormDB.Save(&user).Error
	if err != nil {
		log.Fatalln(err)
	}
	user = User{}
	ormDB.First(&user)
	log.Println(user)

	var users []User
	ormDB.Where("activated_at > '2019-10-31 15:06:00.772'").Find(&users)
	log.Println(users)

}
