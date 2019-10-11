package initialize

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/liov/hoper/go/v1/h_reflect"
	"github.com/liov/hoper/go/v1/initialize/dao"
	"github.com/liov/hoper/go/v1/utils/log"
)

const (
	MYSQL = "mysql"
	POSTGRES = "postgres"
	SQLite = "sqlite3"
)


func (i *Init) P2DB(conf interface{}) {
	dbConf:=DatabaseConfig{}
	if exist := h_reflect.GetExpectTypeValue(conf,&dbConf);!exist{
		return
	}
	var url string
	if dbConf.Type == MYSQL {
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			dbConf.User, dbConf.Password, dbConf.Host,
			dbConf.Port, dbConf.Database, dbConf.Charset)
	} else if dbConf.Type == POSTGRES {
		url = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			dbConf.Host, dbConf.User, dbConf.Database, dbConf.Password)
	} else if dbConf.Type == SQLite {
		url = "../../static/db/hoper.db"
	}
	db, err := gorm.Open(dbConf.Type, url)

	if err != nil {
		log.Error(err)
		os.Exit(10)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return dbConf.TablePrefix + defaultTableName
	}

	if i.Env != PRODUCT {
		//b不set输出空白
		db.SetLogger(stdlog.New(os.Stdout, "", 3))
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(dbConf.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbConf.MaxOpenConns)
	dao.Dao.SetDB(db)
}
