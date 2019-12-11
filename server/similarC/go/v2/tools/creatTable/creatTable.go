package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/fs"
)

var (
	ormDB *gorm.DB
)

func init() {
	path, err := fs.FindFile("config/add-config.toml")
	if err != nil {
		log.Fatal(err)
	}
	var config = struct {
		DataBase struct {
			User     string
			PassWord string
			Host     string
		}
	}{}
	err = configor.New(&configor.Config{Debug: true}).
		Load(&config, path)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", config.DataBase.User, config.DataBase.PassWord, config.DataBase.Host)

	ormDB, err = gorm.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}
	ormDB.SingularTable(true)
	ormDB.LogMode(true)
}

var userMod = []interface{}{
	&model.User{},
	&model.UserExtend{},
	&model.UserActionLog{},
	&model.UserBannedLog{},
	&model.UserFollow{},
	&model.UserScoreLog{},
	&model.UserFollowLog{},
	&model.Education{},
	&model.Work{},
}

func main() {
	ormDB.DropTable(userMod...)
	ormDB.CreateTable(userMod...)
}
