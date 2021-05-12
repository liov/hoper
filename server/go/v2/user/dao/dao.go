package dao

import (
	"database/sql"
	contexti "github.com/liov/hoper/go/v2/tiga/context"
	"net/smtp"

	"github.com/cockroachdb/pebble"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/tiga/initialize"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/utils/dao/db/gorm"
	"github.com/liov/hoper/go/v2/utils/log"
	"gorm.io/gorm"
)

//原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
//其实init主要就是配置文件数据库连接，可以理解为init放进dao
var Dao *dao = &dao{}

type userDao struct {
	ctxi *contexti.Ctx
}

func GetDao(ctx *contexti.Ctx) *userDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &userDao{ctx}
}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   *gorm.DB `config:"database"`
	StdDB    *sql.DB
	PebbleDB *pebble.DB
	// RedisPool Redis连接池
	Redis *redis.Client
	Cache *ristretto.Cache
	//elastic
	MailAuth smtp.Auth `config:"mail"`
}

// CloseDao close the resource.
func (d *dao) Close() {
	if d.PebbleDB != nil {
		d.PebbleDB.Close()
	}
	if d.Redis != nil {
		d.Redis.Close()
	}
	if d.GORMDB != nil {
		rawDB, _ := d.GORMDB.DB()
		rawDB.Close()
	}
}

func (d *dao) Init() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

	d.StdDB, _ = db.DB()
}

func (d *dao) GetDB(log *log.Logger) *gorm.DB {
	if initialize.InitConfig.Env == initialize.DEVELOPMENT {
		return d.GORMDB
	}
	return d.GORMDB.Session(&gorm.Session{
		Logger: &gormi.SQLLogger{Logger: log.Logger,
			Config: &conf.Conf.GORMDB.Gorm.Logger,
		}})
}
