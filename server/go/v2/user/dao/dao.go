package dao

import (
	"database/sql"
	"net/smtp"

	"github.com/cockroachdb/pebble"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

//原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
//其实init主要就是配置文件数据库连接，可以理解为init放进dao
var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   *gorm.DB
	StdDB    *sql.DB
	PebbleDB *pebble.DB
	// RedisPool Redis连接池
	Redis       *redis.Client
	Cache       *ristretto.Cache
	//elastic
	MailAuth smtp.Auth
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

func (d *dao) Custom() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

	d.StdDB, _ = db.DB()
}
