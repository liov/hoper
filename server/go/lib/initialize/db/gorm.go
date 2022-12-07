package db

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	gormi "github.com/liov/hoper/server/go/lib/utils/dao/db/gorm"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/prometheus"
	stdlog "log"
	"os"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	SQLite   = "sqlite3"
)

type DatabaseConfig struct {
	Type, Charset, Database, Schema, TimeZone string
	Host, User, Password                      string
	TimeFormat                                string
	MaxIdleConns, MaxOpenConns                int
	Port                                      int32
	//bug 字段gorm toml不生效
	Gorm       gormi.GORMConfig
	Prometheus bool
}

func (c *DatabaseConfig) Init() {
	c.Gorm.Init()
	if c.TimeFormat == "" {
		c.TimeZone = "Asia/Shanghai"
	}
	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}
	if c.Charset == "" {
		c.Charset = "utf8"
	}
	if c.Type == "" {
		c.Type = POSTGRES
	}
	if c.Port == 0 {
		if c.Type == MYSQL {
			c.Port = 3306
		}
		if c.Type == POSTGRES {
			c.Port = 5432
		}
	}
}

func (conf *DatabaseConfig) Generate(dialector gorm.Dialector) *gorm.DB {

	var db *gorm.DB
	var err error

	// 默认日志
	logger.Default = logger.New(stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), conf.Gorm.Logger)
	dbConfig := &conf.Gorm.Config
	dbConfig.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}

	db, err = gorm.Open(dialector, dbConfig)
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

type DB struct {
	*gorm.DB `init:"entity"`
	Conf     DatabaseConfig `init:"config"`
}

func (db *DB) Table(name string) *gorm.DB {
	gdb := db.DB.Clauses()
	gdb.Statement.TableExpr = &clause.Expr{SQL: gdb.Statement.Quote(name)}
	gdb.Statement.Table = name
	return gdb
}
