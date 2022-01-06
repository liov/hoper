package inject_dao

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	gormi "github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/prometheus"
	stdlog "log"
	"os"
	"runtime"
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
	MaxIdleConns, MaxOpenConns int
	Port                       int32
	//bug 字段gorm toml不生效
	Gorm       gormi.GORMConfig
	Prometheus bool
}

func (conf *DatabaseConfig) generate() *gorm.DB {
	var url string
	var db *gorm.DB
	var err error

	// 默认日志
	logger.Default = logger.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), conf.Gorm.Logger)
	dbConfig := &conf.Gorm.Config
	dbConfig.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
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
	// 日志
	if initialize.InitConfig.Env != initialize.DEVELOPMENT {
		db.Statement.Logger = &gormi.SQLLogger{Logger: log.Default.Logger,
			Config: &conf.Gorm.Logger,
		}
	}

	if conf.Prometheus {
		if conf.Type == MYSQL {
			db.Use(prometheus.New(prometheus.Config{
				DBName:          conf.Database,               // 使用 `DBName` 作为指标 label
				RefreshInterval: 15,                          // 指标刷新频率（默认为 15 秒）
				PushAddr:        "prometheus pusher address", // 如果配置了 `PushAddr`，则推送指标
				MetricsCollector: []prometheus.MetricsCollector{
					&prometheus.MySQL{
						VariableNames: []string{"Threads_running"},
					},
				}, // 用户自定义指标
			}))
		}
	}

	rawDB, _ := db.DB()
	rawDB.SetMaxIdleConns(conf.MaxIdleConns)
	rawDB.SetMaxOpenConns(conf.MaxOpenConns)
	return db
}

func (conf *DatabaseConfig) Generate() interface{} {
	return conf.generate()
}
