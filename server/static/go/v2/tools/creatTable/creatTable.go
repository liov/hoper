package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
		CreatTable struct {
			Type     string
			User     string
			Password string
			Host     string
		}
	}{}
	err = configor.New(&configor.Config{Debug: true}).
		Load(&config, path)
	if err != nil {
		log.Fatal(err)
	}
	var url string
	if config.CreatTable.Type == "mysql" {
		url = fmt.Sprintf("%s:%s@tcp(%s:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", config.CreatTable.User, config.CreatTable.Password, config.CreatTable.Host)
	} else if config.CreatTable.Type == "postgres" {
		url = fmt.Sprintf("host=%s user=%s dbname=test sslmode=disable password=%s",
			config.CreatTable.Host, config.CreatTable.User, config.CreatTable.Password)
	}

	ormDB, err = gorm.Open(config.CreatTable.Type, url)
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
