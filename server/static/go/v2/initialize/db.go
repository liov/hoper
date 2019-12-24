package initialize

import (
	"fmt"
	"runtime"

	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	SQLite   = "sqlite3"
)

type DatabaseConfig struct {
	Type         string
	User         string
	Password     string
	Host         string
	Charset      string
	Database     string
	TimeFormat   string
	TablePrefix  string
	MaxIdleConns int
	MaxOpenConns int
	Port         int32
}

func (conf *DatabaseConfig) Generate() *gorm.DB {
	var url string
	if conf.Type == MYSQL {
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conf.User, conf.Password, conf.Host,
			conf.Port, conf.Database, conf.Charset)
	} else if conf.Type == POSTGRES {
		url = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			conf.Host, conf.User, conf.Database, conf.Password)
	} else if conf.Type == SQLite {
		url = "/data/db/sqlite/" + conf.Database + ".db"
		if runtime.GOOS == "windows" {
			url = ".." + url
		}
	}
	db, err := gorm.Open(conf.Type, url)
	if err != nil {
		log.Fatal(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.TablePrefix + defaultTableName
	}
	return db
}

func (i *Init) P2DB() *gorm.DB {
	conf := &DatabaseConfig{}
	if exist := reflect3.GetFieldValue(i.conf, conf); !exist {
		return nil
	}

	db := conf.Generate()
	if i.Env != PRODUCT {
		//b不set输出空白
		//db.SetLogger(gorm.Logger{stdlog.New(os.Stderr, "", 0)})
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)

	//i.closes = append(i.closes,db.Close)
	//closes = append(closes, func() {log.Info("数据库已关闭")})
	return db
}
