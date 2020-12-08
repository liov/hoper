package initialize

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime"
	"time"

	v2 "github.com/liov/hoper/go/v2/utils/dao/db/gorm/v2"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	SQLite   = "sqlite3"
)

type DatabaseConfig struct {
	Type, Charset, Database    string
	Host, User, Password       string
	TimeFormat                 string
	TablePrefix                string
	MaxIdleConns, MaxOpenConns int
	Port                       int32
	LogLevel                   logger.LogLevel
}

func (conf *DatabaseConfig) Generate() *gorm.DB {
	var url string
	var db *gorm.DB
	var err error
	dbConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.TablePrefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if conf.Type == MYSQL {
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conf.User, conf.Password, conf.Host,
			conf.Port, conf.Database, conf.Charset)
		db, err = gorm.Open(mysql.Open(url), dbConfig)
	} else if conf.Type == POSTGRES {
		url = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			conf.Host, conf.User, conf.Database, conf.Password)
		db, err = gorm.Open(postgres.Open(url), dbConfig)
	} else if conf.Type == SQLite {
		url = "/data/db/sqlite/" + conf.Database + ".db"
		if runtime.GOOS == "windows" {
			url = ".." + url
		}
		db, err = gorm.Open(sqlite.Open(url), dbConfig)
	}
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (init *Init) P2DB() *gorm.DB {
	conf := &DatabaseConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return nil
	}

	db := conf.Generate()

	rawDB, _ := db.DB()
	rawDB.SetMaxIdleConns(conf.MaxIdleConns)
	rawDB.SetMaxOpenConns(conf.MaxOpenConns)
	if init.Env == PRODUCT {
		db.Config.Logger = logger.Default.LogMode(logger.Error)
	} else {
		db.Config.Logger = v2.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), &logger.Config{
			LogLevel:      conf.LogLevel,
			Colorful:      true,
			SlowThreshold: 100 * time.Millisecond,
		})
	}
	//i.closes = append(i.closes,db.CloseDao)
	//closes = append(closes, func() {log.Info("数据库已关闭")})
	return db
}
