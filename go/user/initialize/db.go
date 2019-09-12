package initialize

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/liov/hoper/go/user/internal/config"
	"github.com/liov/hoper/go/user/internal/dao"
	"github.com/liov/hoper/go/utls/log"
)

func (i *Init) P2DB() {
	var dbConf = &config.Config.Database
	var url string
	if dbConf.Type == "mysql" {
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			dbConf.User, dbConf.Password, dbConf.Host,
			dbConf.Port, dbConf.Database, dbConf.Charset)
	} else if dbConf.Type == "postgres" {
		url = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			dbConf.Host, dbConf.User, dbConf.Database, dbConf.Password)
	} else if dbConf.Type == "sqlite3" {
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
	//b不set输出空白
	db.SetLogger(stdlog.New(os.Stdout, "", 3))
	if config.Config.Server.Env != PRODUCT {
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(dbConf.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbConf.MaxOpenConns)
	dao.SetDB(db)
}

func (i *Init) P2Redis() {
	var redisConf = &config.Config.Redis
	url := fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port)
	dao.SetRedis(&redis.Pool{
		MaxIdle:     redisConf.MaxIdle,
		MaxActive:   redisConf.MaxActive,
		IdleTimeout: redisConf.IdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}
			if redisConf.Password != "" {
				if _, err := c.Do("AUTH", redisConf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	})
}
